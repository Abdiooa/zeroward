# zeroward

zeroward is a command-line zero-knowledge encryption command-line program to secure user files at all stages(locally, on transmission, and on rest(while they are stored on the cloud storages)), it meant to and for different cloud storage providers, but for the moment it is just integrated with Yandex Cloud Storage.

##  Storage providers

* Yandex Cloud Storage [:page_facing_up:] (https://cloud.yandex.com/en/services/storage)


## Features

* Secures(Encrypt) users files before uploading them on the cloud.
* Uses two key for encryptions/decryptions
* A KEK (Key-Encryption-Key) to encrypt the DEK(Data-Encryption-Key)
* The DEK is used to encrypt files, then attached with the file as metadata
* The KEK is on the possession of the user,
* The KEK is generate from a passphrase defined by the user in the first use
* the KEK is then stored in a secure place on the user system
* The program gives user a user-friendly interracting with his objects/buckets that he has in the cloud storages
* For the encryption of the user files and the keys, AES-256-gcm algorithm is used to ensure the security of user files and keys, which is a very strong algorithm

## Installation & Get-Starting


License
-------

This is free secure software under the terms of the Apache License.