![Go Workflow Action](https://github.com/Ma-Vin/packages-action-app/actions/workflows/go.yml/badge.svg)

# Package Action Application

Application to determine and delete existing versions of [GitHub packages](https://docs.github.com/en/packages). It
uses [GitHub rest api](https://docs.github.com/en/rest/packages/packages?apiVersion=2022-11-28) to get information of a
package and its versions.

:baby_chick: This repository is just a try out of GoLang and GitHub rest api

:rocket: This application should be put to an GitHub Action
at [Ma-Vin/packages-action](https://github.com/Ma-Vin/packages-action) afterwards.

## Usage

The application can handle versions of the type *&lt;major&gt;.&lt;minor&gt;.&lt;patch&gt;* or
*&lt;major&gt;.&lt;minor&gt;.&lt;patch&gt;-SNAPSHOT*. If *minor* or *patch* are missing they will be handled as zero.

The application has to be configured by environment variables

| Environment Variable   | Required           | Default                  | Description                                                                                                                                            |
|------------------------|--------------------|--------------------------|--------------------------------------------------------------------------------------------------------------------------------------------------------|
| GITHUB_REST_API_URL    |                    | *https://api.github.com* | Protocol and host of the GitHub rest api                                                                                                               |
| GITHUB_ORGANIZATION    |                    |                          | :warning: :construction: Not supported yet                                                                                                             |
| GITHUB_USER            | :heavy_check_mark: |                          | GitHub user who is the owner of the packages                                                                                                           |
| PACKAGE_TYPE           | :heavy_check_mark: |                          | The type of package. At the moment only *maven* is supported (In general there exists *npm, maven, rubygems, docker, nuget, container*)                |
| PACKAGE_NAME           | :heavy_check_mark: |                          | The name of the package whose versions should be deleted                                                                                               |
| VERSION_NAME_TO_DELETE |                    |                          | A concrete version to delete (Independent of *NUMBER_MAJOR_TO_KEEP NUMBER_MINOR_TO_KEEP* and *NUMBER_PATCH_TO_KEEP*)                                   |
| DELETE_SNAPSHOTS       |                    | *false*                  | Indicator whether to delete all snapshots or none (Snapshots are excluded from *NUMBER_MAJOR_TO_KEEP NUMBER_MINOR_TO_KEEP* and *NUMBER_PATCH_TO_KEEP*) |
| NUMBER_MAJOR_TO_KEEP   |                    | keep all                 | Positive number of major versions to keep                                                                                                              |
| NUMBER_MINOR_TO_KEEP   |                    | keep all                 | Positive number of minor versions to keep (within a major version)                                                                                     |
| NUMBER_PATCH_TO_KEEP   |                    | keep all                 | Positive number of patch versions to keep (within a minor version)                                                                                     |
| GITHUB_TOKEN           | :heavy_check_mark: |                          | The access token to use for bearer authentication against GitHub rest api                                                                              |
| DRY_RUN                |                    | *true*                   | Indicator whether to print deletion candidates only or to delete versions/package                                                                      | 

At least one deletion indicator of *VERSION_NAME_TO_DELETE, DELETE_SNAPSHOTS, NUMBER_MAJOR_TO_KEEP NUMBER_MINOR_TO_KEEP*
or
*NUMBER_PATCH_TO_KEEP* must be set

:warning: If there will remain an empty package, the whole package will be deleted instead of its versions :warning:

## Sonarcloud analysis

* [![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=ma-vin_package-action-application&branch=release%2Fv1.1&metric=alert_status)](https://sonarcloud.io/summary/new_code?id=ma-vin_package-action-application&branch=release%2Fv1.1)
* [![Reliability Rating](https://sonarcloud.io/api/project_badges/measure?project=ma-vin_package-action-application&branch=release%2Fv1.1&metric=reliability_rating)](https://sonarcloud.io/summary/new_code?id=ma-vin_package-action-application&branch=release%2Fv1.1) [![Bugs](https://sonarcloud.io/api/project_badges/measure?project=ma-vin_package-action-application&branch=release%2Fv1.1&metric=bugs)](https://sonarcloud.io/summary/new_code?id=ma-vin_package-action-application&branch=release%2Fv1.1)
* [![Security Rating](https://sonarcloud.io/api/project_badges/measure?project=ma-vin_package-action-application&branch=release%2Fv1.1&metric=security_rating)](https://sonarcloud.io/summary/new_code?id=ma-vin_package-action-application&branch=release%2Fv1.1)  [![Vulnerabilities](https://sonarcloud.io/api/project_badges/measure?project=ma-vin_package-action-application&branch=release%2Fv1.1&metric=vulnerabilities)](https://sonarcloud.io/summary/new_code?id=ma-vin_package-action-application&branch=release%2Fv1.1)
* [![Maintainability Rating](https://sonarcloud.io/api/project_badges/measure?project=ma-vin_package-action-application&branch=release%2Fv1.1&metric=sqale_rating)](https://sonarcloud.io/summary/new_code?id=ma-vin_package-action-application&branch=release%2Fv1.1)  [![Technical Debt](https://sonarcloud.io/api/project_badges/measure?project=ma-vin_package-action-application&branch=release%2Fv1.1&metric=sqale_index)](https://sonarcloud.io/summary/new_code?id=ma-vin_package-action-application&branch=release%2Fv1.1)  [![Code Smells](https://sonarcloud.io/api/project_badges/measure?project=ma-vin_package-action-application&branch=release%2Fv1.1&metric=code_smells)](https://sonarcloud.io/summary/new_code?id=ma-vin_package-action-application&branch=release%2Fv1.1)
* [![Coverage](https://sonarcloud.io/api/project_badges/measure?project=ma-vin_package-action-application&branch=release%2Fv1.1&metric=coverage)](https://sonarcloud.io/summary/new_code?id=ma-vin_package-action-application&branch=release%2Fv1.1)
* [![Duplicated Lines (%)](https://sonarcloud.io/api/project_badges/measure?project=ma-vin_package-action-application&branch=release%2Fv1.1&metric=duplicated_lines_density)](https://sonarcloud.io/summary/new_code?id=ma-vin_package-action-application&branch=release%2Fv1.1)  [![Lines of Code](https://sonarcloud.io/api/project_badges/measure?project=ma-vin_package-action-application&branch=release%2Fv1.1&metric=ncloc)](https://sonarcloud.io/summary/new_code?id=ma-vin_package-action-application&branch=release%2Fv1.1)