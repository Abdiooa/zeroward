# zeroward

zeroward is a command-line zero-knowledge encryption command-line program to secure user files at all stages(locally, on transmission, and on rest(while they are stored on the cloud storages)), it meant to and for different cloud storage providers, but for the moment it is just integrated with Yandex Cloud Storage.

##  Storage providers

* Yandex Cloud Storage [:page_facing_up:] (https://cloud.yandex.com/en/services/storage)


## Features

* Secures (Encrypts) user files before uploading them to the cloud.
* Utilizes two keys for encryption/decryption:
  * A KEK (Key-Encryption-Key) to encrypt the DEK (Data-Encryption-Key).
  * The DEK is used to encrypt files and is attached to the file as metadata.
* The KEK is in the possession of the user.
* The KEK is generated from a passphrase defined by the user during the first use.
* The KEK is then stored securely on the user's system.
* Provides a user-friendly interface for interacting with objects/buckets on cloud storage.
* Uses the AES-256-GCM algorithm for file and key encryption.

## Installation & Get-Started
### Prerequisites
NOTE: Don't worry about the program it is secure and simple.
You don't need no prerequisites, just follow the process of installation for you os plateform

### Linux/Darwin(MacOs)
#### Use of Homebrew-tools

Ensure Homebrew is installed for Linux and MacOS:

Install [zeroward](https://github.com/Abdiooa/zeroward/):

```
brew install zeroward
```
Upgrade the zeroward CLI program to the latest version:

```
brew upgrade zeroward
```

#### Use of the released packages
Download Released Packages (Linux/MacOS)

Downloading and Installing (Debian)
Download the latest release for Linux (amd64) from the [releases](https://github.com/Abdiooa/zeroward/releases) page:

```
wget https://github.com/Abdiooa/zeroward/releases/vtag/download/zeroward_vtag_linux_amd64.deb
```
please make sure to replace the vtag the tag version you want, you can see the latest one here [releases](https://github.com/Abdiooa/zeroward/releases)
Install the downloaded Debian package using dpkg:
```
sudo dpkg -i zeroward_linux_amd64.deb
```
To uninstall it use this command:
```
sudo apt-get remove zeroward
```
Alternatively, for other architectures or package formats:
```
wget https://github.com/Abdiooa/zeroward/releases/download/vX.Y.Z/zeroward_X.Y.Z_Linux_amd64.tar.gz
tar -zxvf zeroward_X.Y.Z_Linux_amd64.tar.gz
sudo mv zeroward /usr/local/bin/
```
Repeat the process for other architectures or package formats.

#### Windows
1. Download and Install
Visit the releases page on GitHub.

2. Download the latest release zip file (e.g., zeroward_windows_amd64.zip).

3. Extract the contents of the zip file.

4. Open a command prompt in the extracted folder.

5. Run the following command to install zeroward:
```
zeroward.exe install
```
### Usage
```
zeroward --help

# List all buckets
zeroward buckets --accessKeyID accesskeyid --secretAccessKey secretacccesskey

# List objects in a bucket
zeroward objects --bcktname bucketname --accessKeyID accesskeyid --secretAccessKey secretacccesskey

# Upload a file to the cloud
zeroward upload --bcktname bucketname --filePath pathtothefile --passphrase passphrase --objectkey paththefilestoredgonnastored 

# Download a file from the cloud
zeroward download --filePath pathtothefile --objectkey --accessKeyID accesskeyid --secretAccessKey secretacccesskeypaththefilestoredonthecloud --bcktname bucketname --removeAfterDownload y --accessKeyID accesskeyid --secretAccessKey secretacccesskey

```
#### the access id key, the secret access key and the passphrase are just meant to be defined for the first use of the application

License
-------
This is free secure software under the terms of the Apache License.