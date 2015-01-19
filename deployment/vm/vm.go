package vm

import (
	"time"

	bosherr "github.com/cloudfoundry/bosh-agent/errors"
	boshlog "github.com/cloudfoundry/bosh-agent/logger"
	boshsys "github.com/cloudfoundry/bosh-agent/system"
	boshtime "github.com/cloudfoundry/bosh-agent/time"

	bmcloud "github.com/cloudfoundry/bosh-micro-cli/cloud"
	bmconfig "github.com/cloudfoundry/bosh-micro-cli/config"
	bmagentclient "github.com/cloudfoundry/bosh-micro-cli/deployment/agentclient"
	bmas "github.com/cloudfoundry/bosh-micro-cli/deployment/applyspec"
	bmdisk "github.com/cloudfoundry/bosh-micro-cli/deployment/disk"
	bmdeplmanifest "github.com/cloudfoundry/bosh-micro-cli/deployment/manifest"
	bmretrystrategy "github.com/cloudfoundry/bosh-micro-cli/deployment/retrystrategy"
	bmstemcell "github.com/cloudfoundry/bosh-micro-cli/deployment/stemcell"
	bmeventlog "github.com/cloudfoundry/bosh-micro-cli/eventlogger"
)

type VM interface {
	CID() string
	Exists() (bool, error)
	WaitUntilReady(timeout time.Duration, delay time.Duration) error
	Apply(bmstemcell.ApplySpec, bmdeplmanifest.Manifest) error
	UpdateDisks(bmdeplmanifest.DiskPool, bmeventlog.Stage) ([]bmdisk.Disk, error)
	Start() error
	WaitToBeRunning(maxAttempts int, delay time.Duration) error
	AttachDisk(bmdisk.Disk) error
	DetachDisk(bmdisk.Disk) error
	Stop() error
	Disks() ([]bmdisk.Disk, error)
	UnmountDisk(bmdisk.Disk) error
	MigrateDisk() error
	Delete() error
}

type vm struct {
	cid                    string
	vmRepo                 bmconfig.VMRepo
	stemcellRepo           bmconfig.StemcellRepo
	diskDeployer           DiskDeployer
	agentClient            bmagentclient.AgentClient
	cloud                  bmcloud.Cloud
	templatesSpecGenerator bmas.TemplatesSpecGenerator
	applySpecFactory       bmas.Factory
	mbusURL                string
	fs                     boshsys.FileSystem
	logger                 boshlog.Logger
	logTag                 string
}

func NewVM(
	cid string,
	vmRepo bmconfig.VMRepo,
	stemcellRepo bmconfig.StemcellRepo,
	diskDeployer DiskDeployer,
	agentClient bmagentclient.AgentClient,
	cloud bmcloud.Cloud,
	templatesSpecGenerator bmas.TemplatesSpecGenerator,
	applySpecFactory bmas.Factory,
	mbusURL string,
	fs boshsys.FileSystem,
	logger boshlog.Logger,
) VM {
	return &vm{
		cid:          cid,
		vmRepo:       vmRepo,
		stemcellRepo: stemcellRepo,
		diskDeployer: diskDeployer,
		agentClient:  agentClient,
		cloud:        cloud,
		templatesSpecGenerator: templatesSpecGenerator,
		applySpecFactory:       applySpecFactory,
		mbusURL:                mbusURL,
		fs:                     fs,
		logger:                 logger,
		logTag:                 "vm",
	}
}

func (vm *vm) CID() string {
	return vm.cid
}

func (vm *vm) Exists() (bool, error) {
	exists, err := vm.cloud.HasVM(vm.cid)
	if err != nil {
		return false, bosherr.WrapErrorf(err, "Checking existance of VM '%s'", vm.cid)
	}
	return exists, nil
}

func (vm *vm) WaitUntilReady(timeout time.Duration, delay time.Duration) error {
	agentPingRetryable := bmagentclient.NewPingRetryable(vm.agentClient)
	timeService := boshtime.NewConcreteService()
	agentPingRetryStrategy := bmretrystrategy.NewTimeoutRetryStrategy(timeout, delay, agentPingRetryable, timeService, vm.logger)
	return agentPingRetryStrategy.Try()
}

func (vm *vm) Apply(stemcellApplySpec bmstemcell.ApplySpec, deploymentManifest bmdeplmanifest.Manifest) error {
	vm.logger.Debug(vm.logTag, "Stopping agent")

	err := vm.agentClient.Stop()
	if err != nil {
		return bosherr.WrapError(err, "Stopping agent")
	}

	vm.logger.Debug(vm.logTag, "Rendering job templates")
	renderedJobDir, err := vm.fs.TempDir("instance-updater-render-job")
	if err != nil {
		return bosherr.WrapError(err, "Creating rendered job directory")
	}
	defer vm.fs.RemoveAll(renderedJobDir)

	deploymentJob := deploymentManifest.Jobs[0]
	jobProperties, err := deploymentJob.Properties()
	if err != nil {
		return bosherr.WrapError(err, "Stringifying job properties")
	}

	networksSpec, err := deploymentManifest.NetworksSpec(deploymentJob.Name)
	if err != nil {
		return bosherr.WrapError(err, "Stringifying job properties")
	}

	templatesSpec, err := vm.templatesSpecGenerator.Create(
		deploymentJob,
		stemcellApplySpec.Job,
		deploymentManifest.Name,
		jobProperties,
		vm.mbusURL,
	)
	if err != nil {
		return bosherr.WrapError(err, "Generating templates spec")
	}

	vm.logger.Debug(vm.logTag, "Creating apply spec")
	agentApplySpec := vm.applySpecFactory.Create(
		stemcellApplySpec,
		deploymentManifest.Name,
		deploymentJob.Name,
		networksSpec,
		templatesSpec.BlobID,
		templatesSpec.ArchiveSha1,
		templatesSpec.ConfigurationHash,
	)

	vm.logger.Debug(vm.logTag, "Sending apply message to the agent with '%#v'", agentApplySpec)
	err = vm.agentClient.Apply(agentApplySpec)
	if err != nil {
		return bosherr.WrapError(err, "Sending apply spec to agent")
	}

	return nil
}

func (vm *vm) UpdateDisks(diskPool bmdeplmanifest.DiskPool, eventLoggerStage bmeventlog.Stage) ([]bmdisk.Disk, error) {
	disks, err := vm.diskDeployer.Deploy(diskPool, vm.cloud, vm, eventLoggerStage)
	if err != nil {
		return disks, bosherr.WrapError(err, "Deploying disk")
	}
	return disks, nil
}

func (vm *vm) Start() error {
	return vm.agentClient.Start()
}

func (vm *vm) WaitToBeRunning(maxAttempts int, delay time.Duration) error {
	agentGetStateRetryable := bmagentclient.NewGetStateRetryable(vm.agentClient)
	agentGetStateRetryStrategy := bmretrystrategy.NewAttemptRetryStrategy(maxAttempts, delay, agentGetStateRetryable, vm.logger)
	return agentGetStateRetryStrategy.Try()
}

func (vm *vm) AttachDisk(disk bmdisk.Disk) error {
	err := vm.cloud.AttachDisk(vm.cid, disk.CID())
	if err != nil {
		return bosherr.WrapError(err, "Attaching disk in the cloud")
	}

	err = vm.agentClient.MountDisk(disk.CID())
	if err != nil {
		return bosherr.WrapError(err, "Mounting disk")
	}

	return nil
}

func (vm *vm) DetachDisk(disk bmdisk.Disk) error {
	err := vm.cloud.DetachDisk(vm.cid, disk.CID())
	if err != nil {
		return bosherr.WrapError(err, "Detaching disk in the cloud")
	}

	return nil
}

func (vm *vm) Stop() error {
	return vm.agentClient.Stop()
}

func (vm *vm) Disks() ([]bmdisk.Disk, error) {
	result := []bmdisk.Disk{}

	disks, err := vm.agentClient.ListDisk()
	if err != nil {
		return result, bosherr.WrapError(err, "Listing vm disks")
	}

	for _, diskCID := range disks {
		disk := bmdisk.NewDisk(bmconfig.DiskRecord{CID: diskCID}, nil, nil)
		result = append(result, disk)
	}

	return result, nil
}

func (vm *vm) UnmountDisk(disk bmdisk.Disk) error {
	return vm.agentClient.UnmountDisk(disk.CID())
}

func (vm *vm) MigrateDisk() error {
	return vm.agentClient.MigrateDisk()
}

func (vm *vm) Delete() error {
	deleteErr := vm.cloud.DeleteVM(vm.cid)
	if deleteErr != nil {
		// allow VMNotFoundError for idempotency
		cloudErr, ok := deleteErr.(bmcloud.Error)
		if !ok || cloudErr.Type() != bmcloud.VMNotFoundError {
			return bosherr.WrapError(deleteErr, "Deleting vm in the cloud")
		}
	}

	err := vm.vmRepo.ClearCurrent()
	if err != nil {
		return bosherr.WrapError(err, "Deleting vm from vm repo")
	}

	err = vm.stemcellRepo.ClearCurrent()
	if err != nil {
		return bosherr.WrapError(err, "Clearing current stemcell from stemcell repo")
	}

	// returns bmcloud.Error only if it is a VMNotFoundError
	return deleteErr
}