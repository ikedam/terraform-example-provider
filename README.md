# Terraform のプライベートレジストリーのデモ

## 動作テスト

```
docker compose run --rm terraform init
docker compose run --rm terraform plan
```

## バイナリーを自分で更新する場合の手順

### 事前準備

gpg キーを作成しておくこと

```
gpg --full-generate-key
```

### プライベートレジストリーの構築手順

```
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o terraform-provider-hello .
mv terraform-provider-hello registry/dl/terraform-providers/example/hello/1.0.0/download/linux/amd64/terraform-provider-hello_1.0.0_linux_amd64
cd registry/dl/terraform-providers/example/hello/1.0.0/download/linux/amd64
zip terraform-provider-hello_1.0.0_linux_amd64.zip terraform-provider-hello_1.0.0_linux_amd64
sha256sum terraform-provider-hello_1.0.0_linux_amd64.zip >terraform-provider-hello_1.0.0_linux_amd64_SHA256SUMS
gpg --output terraform-provider-hello_1.0.0_linux_amd64_SHA256SUMS.sig --detach-sig terraform-provider-hello_1.0.0_linux_amd64_SHA256SUMS
```

index.json 用の情報の収集:

```
gpg --list-secret-keys --keyid-format=long
```

```
gpg --armor --export キーID | jq -Rs -c . 
```


## メモ

curl https://registry.terraform.io/.well-known/terraform.json
{"modules.v1":"/v1/modules/","providers.v1":"/v1/providers/"}

curl https://registry.terraform.io/v1/providers/hashicorp/aws/versions

```
{
  "id": "hashicorp/aws",
  "versions": [
    ...
    {
      "version": "5.93.0",
      "protocols": [
        "5.0"
      ],
      "platforms": [
        {
          "os": "openbsd",
          "arch": "arm"
        },
        {
          "os": "linux",
          "arch": "arm64"
        },
        {
          "os": "linux",
          "arch": "amd64"
        },
        ...
```

curl https://registry.terraform.io/v1/providers/hashicorp/aws/5.93.0/download/linux/amd64

```
{
  "protocols": [
    "5.0"
  ],
  "os": "linux",
  "arch": "amd64",
  "filename": "terraform-provider-aws_5.93.0_linux_amd64.zip",
  "download_url": "https://releases.hashicorp.com/terraform-provider-aws/5.93.0/terraform-provider-aws_5.93.0_linux_amd64.zip",
  "shasums_url": "https://releases.hashicorp.com/terraform-provider-aws/5.93.0/terraform-provider-aws_5.93.0_SHA256SUMS",
  "shasums_signature_url": "https://releases.hashicorp.com/terraform-provider-aws/5.93.0/terraform-provider-aws_5.93.0_SHA256SUMS.72D7468F.sig",
  "shasum": "41cf69a525f0fbe0fdb71d26be7ff5e20bb90ccdf5af32c83ed53f0ca2f071b5",
  "signing_keys": {
    "gpg_public_keys": [
      {
        "key_id": "34365D9472D7468F",
        "ascii_armor": "(略)",
        "trust_signature": "",
        "source": "HashiCorp",
        "source_url": "https://www.hashicorp.com/security.html"
      }
    ]
  }
}
```

curl https://releases.hashicorp.com/terraform-provider-aws/5.93.0/terraform-provider-aws_5.93.0_SHA256SUMS
```
43055bdd0786855cf7242638a74b579f74f4f1a8e7c7e5e0e50230c8f6b908cb  terraform-provider-aws_5.93.0_darwin_amd64.zip
...
41cf69a525f0fbe0fdb71d26be7ff5e20bb90ccdf5af32c83ed53f0ca2f071b5  terraform-provider-aws_5.93.0_linux_amd64.zip
...
```


$ curl https://releases.hashicorp.com/terraform-provider-aws/5.93.0/terraform-provider-aws_5.93.0_SHA256SUMS.72D7468F.sig
Warning: Binary output can mess up your terminal. Use "--output -" to tell
Warning: curl to output it to your terminal anyway, or consider "--output
Warning: <FILE>" to save to a file.


$ curl -O https://releases.hashicorp.com/terraform-provider-aws/5.93.0/terraform-provider-aws_5.93.0_linux_amd64.zip
$ zipinfo terraform-provider-aws_5.93.0_linux_amd64.zip
Archive:  terraform-provider-aws_5.93.0_linux_amd64.zip
Zip file size: 150657074 bytes, number of entries: 2
-rw-r--r--  2.0 unx    16761 bX defN 25-Mar-28 03:37 LICENSE.txt
-rwxr-xr-x  2.0 unx 690688152 bX defN 25-Mar-28 03:05 terraform-provider-aws_v5.93.0_x5
2 files, 690704913 bytes uncompressed, 150656744 bytes compressed:  78.2%






https://developer.hashicorp.com/terraform/internals/provider-registry-protocol#provider-addresses

プロバイダー名を決める

(ドメイン名)/(ネームスペース)/(タイプ) みたいな感じです。

今回は ikedam.jp/example/hello にしました。


サービスディスカバリーの設定
https://developer.hashicorp.com/terraform/internals/provider-registry-protocol#service-discovery

/.well-known/terraform.json
{
  "providers.v1": "/dl/terraform-providers/
}


Content-Type: application/json

でないといけないのがちょっとハードル。


バージョン一覧ファイル
https://developer.hashicorp.com/terraform/internals/provider-registry-protocol#list-available-versions


/dl/terraform-providers/example/hello/versions

* 「/dl/terraform-providers/」の部分はサービスディスカバリーに対応付ける。
* protocols という項目があるが、使用しているプロトコルバージョンを特定する方法が不明。



```
{
  "versions": [
    {
      "version": "1.0.0",
      "platforms": [
        {"os": "linux", "arch": "amd64"}
      ]
    }
  ]
}
```

プロバイダーパッケージ情報
https://developer.hashicorp.com/terraform/internals/provider-registry-protocol#find-a-provider-package


/dl/terraform-providers/example/hello/1.0.0/download/linux/amd64

* 「/dl/terraform-providers/」の部分はサービスディスカバリーに対応付ける。
* protocols という項目があるが、使用しているプロトコルバージョンを特定する方法が不明。



```
{
  "versions": [
    {
      "version": "1.0.0",
      "platforms": [
        {"os": "linux", "arch": "amd64"}
      ]
    }
  ]
}
```



```
zip terraform-provider-hello_1.0.0_linux_amd64.zip terraform-provider-hello_1.0.0_linux_amd64rld-provider
sha256sum terraform-provider-hello_1.0.0_linux_amd64.zip  >terraform-provider-hello_1.0.0_linux_amd64_SHA256SUMS
```



GPG キーの作成

```
$ gpg --full-generate-key
gpg (GnuPG) 2.4.7-unknown; Copyright (C) 2024 g10 Code GmbH
This is free software: you are free to change and redistribute it.
There is NO WARRANTY, to the extent permitted by law.

Please select what kind of key you want:
   (1) RSA and RSA
   (2) DSA and Elgamal
   (3) DSA (sign only)
   (4) RSA (sign only)
   (9) ECC (sign and encrypt) *default*
  (10) ECC (sign only)
  (14) Existing key from card
Your selection?
Please select which elliptic curve you want:
   (1) Curve 25519 *default*
   (4) NIST P-384
   (6) Brainpool P-256
Your selection?
Please specify how long the key should be valid.
         0 = key does not expire
      <n>  = key expires in n days
      <n>w = key expires in n weeks
      <n>m = key expires in n months
      <n>y = key expires in n years
Key is valid for? (0)
Key does not expire at all
Is this correct? (y/N) y

GnuPG needs to construct a user ID to identify your key.

Real name: ikedam
Email address: ikedam@example.com
Comment:
You selected this USER-ID:
    "ikedam <ikedam@example.com>"

Change (N)ame, (C)omment, (E)mail or (O)kay/(Q)uit? O
We need to generate a lot of random bytes. It is a good idea to perform
some other action (type on the keyboard, move the mouse, utilize the
disks) during the prime generation; this gives the random number
generator a better chance to gain enough entropy.
We need to generate a lot of random bytes. It is a good idea to perform
some other action (type on the keyboard, move the mouse, utilize the
disks) during the prime generation; this gives the random number
generator a better chance to gain enough entropy.
gpg: /home/yasuke/.gnupg/trustdb.gpg: trustdb created
gpg: directory '/home/yasuke/.gnupg/openpgp-revocs.d' created
gpg: revocation certificate stored as '/home/yasuke/.gnupg/openpgp-revocs.d/89522820B01D816002BA019402049A7C5E3F87D5.rev'
public and secret key created and signed.

pub   ed25519 2025-03-30 [SC]
      89522820B01D816002BA019402049A7C5E3F87D5
uid                      ikedam <ikedam@example.com>
sub   cv25519 2025-03-30 [E]
```


$ gpg --output terraform-provider-hello_1.0.0_linux_amd64_SHA256SUMS.sig --detach-sig terraform-provider-hello_1.0.0_linux_amd64_SHA256SUMS

$ gpg --list-secret-keys --keyid-format=long
gpg: checking the trustdb
gpg: marginals needed: 3  completes needed: 1  trust model: pgp
gpg: depth: 0  valid:   1  signed:   0  trust: 0-, 0q, 0n, 0m, 0f, 1u
[keyboxd]
---------
sec   ed25519/02049A7C5E3F87D5 2025-03-30 [SC]
      89522820B01D816002BA019402049A7C5E3F87D5
uid                 [ultimate] ikedam <ikedam@example.com>
ssb   cv25519/E9F005302510266D 2025-03-30 [E]


$ gpg --armor --export 02049A7C5E3F87D5
-----BEGIN PGP PUBLIC KEY BLOCK-----

mDMEZ+jtJRYJKwYBBAHaRw8BAQdACZocwNh8Qc9Cty2owoQRocNeYP2GOCIgEaeE
ZIHXoT60G2lrZWRhbSA8aWtlZGFtQGV4YW1wbGUuY29tPoiTBBMWCgA7FiEEiVIo
ILAdgWACugGUAgSafF4/h9UFAmfo7SUCGwMFCwkIBwICIgIGFQoJCAsCBBYCAwEC
HgcCF4AACgkQAgSafF4/h9WkdAD9HtcYcuS0Tj5ZHP69tpDhsYxBFm1ekHNbYxVh
56OdVP0A/ib6v4mnaOpUx3FfbqaRItnihKZIhOkBHm9qMzym9QkOuDgEZ+jtJRIK
KwYBBAGXVQEFAQEHQGYhlsH43WcQq741wNr9r3sfTAIazr0puPRaSlNG6KNmAwEI
B4h4BBgWCgAgFiEEiVIoILAdgWACugGUAgSafF4/h9UFAmfo7SUCGwwACgkQAgSa
fF4/h9UK1wD+JFQfPLXud/jXPbGUNc0UAbyi517ZGq4PnfxKXSN/YrYBANtgaVnX
x2wYaztOim7KMnYaG8PC0/eaZadG0gS68UAL
=c1st
-----END PGP PUBLIC KEY BLOCK-----


$ gpg --armor --export 02049A7C5E3F87D5 | jq -Rs -c .
"-----BEGIN PGP PUBLIC KEY BLOCK-----\n\nmDMEZ+jtJRYJKwYBBAHaRw8BAQdACZocwNh8Qc9Cty2owoQRocNeYP2GOCIgEaeE\nZIHXoT60G2lrZWRhbSA8aWtlZGFtQGV4YW1wbGUuY29tPoiTBBMWCgA7FiEEiVIo\nILAdgWACugGUAgSafF4/h9UFAmfo7SUCGwMFCwkIBwICIgIGFQoJCAsCBBYCAwEC\nHgcCF4AACgkQAgSafF4/h9WkdAD9HtcYcuS0Tj5ZHP69tpDhsYxBFm1ekHNbYxVh\n56OdVP0A/ib6v4mnaOpUx3FfbqaRItnihKZIhOkBHm9qMzym9QkOuDgEZ+jtJRIK\nKwYBBAGXVQEFAQEHQGYhlsH43WcQq741wNr9r3sfTAIazr0puPRaSlNG6KNmAwEI\nB4h4BBgWCgAgFiEEiVIoILAdgWACugGUAgSafF4/h9UFAmfo7SUCGwwACgkQAgSa\nfF4/h9UK1wD+JFQfPLXud/jXPbGUNc0UAbyi517ZGq4PnfxKXSN/YrYBANtgaVnX\nx2wYaztOim7KMnYaG8PC0/eaZadG0gS68UAL\n=c1st\n-----END PGP PUBLIC KEY BLOCK-----\n"

