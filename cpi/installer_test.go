package cpi_test

import (
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/cloudfoundry/bosh-micro-cli/cpi"

	boshlog "github.com/cloudfoundry/bosh-agent/logger"

	bmcpiinstall "github.com/cloudfoundry/bosh-micro-cli/cpi/install"
	bmdepl "github.com/cloudfoundry/bosh-micro-cli/deployment"
	bmrel "github.com/cloudfoundry/bosh-micro-cli/release"

	fakesys "github.com/cloudfoundry/bosh-agent/system/fakes"
	fakebmcloud "github.com/cloudfoundry/bosh-micro-cli/cloud/fakes"
	fakebmcomp "github.com/cloudfoundry/bosh-micro-cli/cpi/compile/fakes"
	fakebmjobi "github.com/cloudfoundry/bosh-micro-cli/cpi/install/fakes"
	fakebmrel "github.com/cloudfoundry/bosh-micro-cli/release/fakes"
	testfakes "github.com/cloudfoundry/bosh-micro-cli/testutils/fakes"
	fakebmui "github.com/cloudfoundry/bosh-micro-cli/ui/fakes"
)

var _ = Describe("Installer", func() {

	var (
		fakeFs               *fakesys.FakeFileSystem
		fakeExtractor        *testfakes.FakeMultiResponseExtractor
		fakeReleaseValidator *fakebmrel.FakeValidator
		fakeReleaseCompiler  *fakebmcomp.FakeReleaseCompiler
		fakeJobInstaller     *fakebmjobi.FakeJobInstaller
		fakeCloudFactory     *fakebmcloud.FakeFactory
		fakeUI               *fakebmui.FakeUI

		deploymentManifestPath string
		cpiInstaller           Installer
	)
	BeforeEach(func() {
		fakeFs = fakesys.NewFakeFileSystem()
		fakeExtractor = testfakes.NewFakeMultiResponseExtractor()
		fakeReleaseValidator = fakebmrel.NewFakeValidator()
		fakeReleaseCompiler = fakebmcomp.NewFakeReleaseCompiler()
		fakeJobInstaller = fakebmjobi.NewFakeJobInstaller()
		fakeCloudFactory = fakebmcloud.NewFakeFactory()
		fakeUI = &fakebmui.FakeUI{}
		logger := boshlog.NewLogger(boshlog.LevelNone)

		deploymentManifestPath = "/fake/manifest.yml"
		cpiInstaller = NewInstaller(fakeUI, fakeFs, fakeExtractor, fakeReleaseValidator, fakeReleaseCompiler, fakeJobInstaller, fakeCloudFactory, logger)
	})

	Describe("Extract", func() {
		var (
			releaseTarballPath string
		)
		BeforeEach(func() {
			releaseTarballPath = "/fake/release.tgz"
			fakeFs.WriteFileString(releaseTarballPath, "")
		})

		Context("when a extracted release directory can be created", func() {
			var (
				release    bmrel.Release
				releaseJob bmrel.Job
			)

			BeforeEach(func() {
				fakeFs.TempDirDir = "/extracted-release-path"

				releasePackage := &bmrel.Package{
					Name:          "fake-release-package-name",
					Fingerprint:   "fake-release-package-fingerprint",
					SHA1:          "fake-release-package-sha1",
					Dependencies:  []*bmrel.Package{},
					ExtractedPath: "/extracted-release-path/extracted_packages/fake-release-package-name",
				}

				releaseJob = bmrel.Job{
					Name:          "fake-release-job-name",
					Fingerprint:   "fake-release-job-fingerprint",
					SHA1:          "fake-release-job-sha1",
					ExtractedPath: "/extracted-release-path/extracted_jobs/fake-release-job-name",
					Templates: map[string]string{
						"cpi.erb":     "bin/cpi",
						"cpi.yml.erb": "config/cpi.yml",
					},
					PackageNames: []string{releasePackage.Name},
					Packages:     []*bmrel.Package{releasePackage},
					Properties:   map[string]bmrel.PropertyDefinition{},
				}

				fakeFS := fakesys.NewFakeFileSystem()

				release = bmrel.NewRelease(
					"fake-release-name",
					"fake-release-version",
					[]bmrel.Job{releaseJob},
					[]*bmrel.Package{releasePackage},
					"/extracted-release-path",
					fakeFS,
				)

				releaseContents := `---
name: fake-release-name
version: fake-release-version

packages:
- name: fake-release-package-name
  version: fake-release-package-version
  fingerprint: fake-release-package-fingerprint
  sha1: fake-release-package-sha1
  dependencies: []
jobs:
- name: fake-release-job-name
  version: fake-release-job-version
  fingerprint: fake-release-job-fingerprint
  sha1: fake-release-job-sha1
`
				fakeFs.WriteFileString("/extracted-release-path/release.MF", releaseContents)
				jobManifestContents := `---
name: fake-release-job-name
templates:
 cpi.erb: bin/cpi
 cpi.yml.erb: config/cpi.yml

packages:
 - fake-release-package-name

properties: {}
`
				fakeFs.WriteFileString("/extracted-release-path/extracted_jobs/fake-release-job-name/job.MF", jobManifestContents)
			})

			Context("and the tarball is a valid BOSH release", func() {
				It("extracts the release to the ExtractedPath", func() {
					release, err := cpiInstaller.Extract(releaseTarballPath)
					Expect(err).NotTo(HaveOccurred())

					expectedPackage := &bmrel.Package{
						Name:          "fake-release-package-name",
						Fingerprint:   "fake-release-package-fingerprint",
						SHA1:          "fake-release-package-sha1",
						ExtractedPath: "/extracted-release-path/extracted_packages/fake-release-package-name",
						Dependencies:  []*bmrel.Package{},
					}
					expectedRelease := bmrel.NewRelease(
						"fake-release-name",
						"fake-release-version",
						[]bmrel.Job{
							{
								Name:          "fake-release-job-name",
								Fingerprint:   "fake-release-job-fingerprint",
								SHA1:          "fake-release-job-sha1",
								ExtractedPath: "/extracted-release-path/extracted_jobs/fake-release-job-name",
								Templates: map[string]string{
									"cpi.erb":     "bin/cpi",
									"cpi.yml.erb": "config/cpi.yml",
								},
								PackageNames: []string{
									"fake-release-package-name",
								},
								Packages:   []*bmrel.Package{expectedPackage},
								Properties: map[string]bmrel.PropertyDefinition{},
							},
						},
						[]*bmrel.Package{expectedPackage},
						"/extracted-release-path",
						fakeFs,
					)

					Expect(release).To(Equal(expectedRelease))

					Expect(fakeFs.FileExists("/extracted-release-path")).To(BeTrue())
					Expect(fakeFs.FileExists("/extracted-release-path/extracted_packages/fake-release-package-name")).To(BeTrue())
					Expect(fakeFs.FileExists("/extracted-release-path/extracted_jobs/fake-release-job-name")).To(BeTrue())
				})
			})

			Context("and the tarball is not a valid BOSH release", func() {
				BeforeEach(func() {
					fakeFs.WriteFileString("/extracted-release-path/release.MF", `{}`)
					fakeReleaseValidator.ValidateError = errors.New("fake-error")
				})

				It("returns an error", func() {
					_, err := cpiInstaller.Extract(releaseTarballPath)
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(ContainSubstring("fake-error"))
				})
			})

			Context("and the tarball cannot be read", func() {
				It("returns an error", func() {
					fakeExtractor.SetDecompressBehavior(releaseTarballPath, "/extracted-release-path", errors.New("fake-error"))
					_, err := cpiInstaller.Extract(releaseTarballPath)
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(ContainSubstring("Reading CPI release from `/fake/release.tgz'"))
					Expect(fakeUI.Errors).To(ContainElement("CPI release at `/fake/release.tgz' is not a BOSH release"))
				})
			})
		})

		Context("when a extracted release path cannot be created", func() {
			BeforeEach(func() {
				fakeFs.TempDirError = errors.New("fake-tmp-dir-error")
			})

			It("returns an error", func() {
				_, err := cpiInstaller.Extract(releaseTarballPath)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("fake-tmp-dir-error"))
				Expect(err.Error()).To(ContainSubstring("Creating temp directory"))
				Expect(fakeUI.Errors).To(ContainElement("Could not create a temporary directory"))
			})
		})
	})

	Describe("Install", func() {
		var (
			deployment bmdepl.Deployment
			release    bmrel.Release
			releaseJob bmrel.Job
		)
		BeforeEach(func() {
			fakeFs.WriteFileString(deploymentManifestPath, "")

			deployment = bmdepl.Deployment{}

			releasePackage := &bmrel.Package{
				Name:          "fake-release-package-name",
				Fingerprint:   "fake-release-package-fingerprint",
				SHA1:          "fake-release-package-sha1",
				Dependencies:  []*bmrel.Package{},
				ExtractedPath: "/extracted-release-path/extracted_packages/fake-release-package-name",
			}

			releaseJob = bmrel.Job{
				Name:          "fake-release-job-name",
				Fingerprint:   "fake-release-job-fingerprint",
				SHA1:          "fake-release-job-sha1",
				ExtractedPath: "/extracted-release-path/extracted_jobs/fake-release-job-name",
				Templates: map[string]string{
					"cpi.erb":     "bin/cpi",
					"cpi.yml.erb": "config/cpi.yml",
				},
				PackageNames: []string{releasePackage.Name},
				Packages:     []*bmrel.Package{releasePackage},
				Properties:   map[string]bmrel.PropertyDefinition{},
			}

			fakeFS := fakesys.NewFakeFileSystem()

			release = bmrel.NewRelease(
				"fake-release-name",
				"fake-release-version",
				[]bmrel.Job{releaseJob},
				[]*bmrel.Package{releasePackage},
				"/extracted-release-path",
				fakeFS,
			)
		})

		Context("and the tarball is a valid BOSH release", func() {
			var (
				installedJob   bmcpiinstall.InstalledJob
				installedJobs  []bmcpiinstall.InstalledJob
				installedCloud *fakebmcloud.FakeCloud
			)

			BeforeEach(func() {
				deployment = bmdepl.Deployment{
					Name:          "fake-deployment-name",
					RawProperties: map[interface{}]interface{}{},
					Jobs: []bmdepl.Job{
						bmdepl.Job{
							Name:      "fake-deployment-job-name",
							Instances: 1,
							Templates: []bmdepl.ReleaseJobRef{
								bmdepl.ReleaseJobRef{
									Name:    "fake-release-job-name",
									Release: "fake-release-name",
								},
							},
						},
					},
				}

				fakeReleaseCompiler.SetCompileBehavior(release, deployment, nil)

				installedJob = bmcpiinstall.InstalledJob{
					Name: "fake-release-job-name",
					Path: "/extracted-release-path/fake-release-job-name",
				}
				fakeJobInstaller.SetInstallBehavior(releaseJob, installedJob, nil)

				installedJobs = []bmcpiinstall.InstalledJob{installedJob}
				installedCloud = fakebmcloud.NewFakeCloud()
				fakeCloudFactory.SetNewCloudBehavior(installedJobs, installedCloud, nil)
			})

			It("compiles the release", func() {
				_, err := cpiInstaller.Install(deployment, release)
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeReleaseCompiler.CompileInputs[0].Deployment).To(Equal(deployment))
			})

			It("installs the deployment jobs", func() {
				_, err := cpiInstaller.Install(deployment, release)
				Expect(err).NotTo(HaveOccurred())

				Expect(fakeJobInstaller.JobInstallInputs).To(Equal(
					[]fakebmjobi.JobInstallInput{
						fakebmjobi.JobInstallInput{
							Job: releaseJob,
						},
					},
				))
			})

			It("returns a cloud wrapper around the installed CPI", func() {
				cloud, err := cpiInstaller.Install(deployment, release)
				Expect(err).NotTo(HaveOccurred())

				Expect(cloud).To(Equal(installedCloud))
			})
		})

		Context("when compilation fails", func() {
			It("returns an error", func() {
				fakeReleaseCompiler.SetCompileBehavior(release, deployment, errors.New("fake-compile-error"))

				_, err := cpiInstaller.Install(deployment, release)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("fake-compile-error"))
				Expect(fakeUI.Errors).To(ContainElement("Could not compile CPI release"))
			})
		})
	})
})