# Bindplane Terraform Production Example

## Description

This basic example will:
- deploy a collector to google cloud and azure
- deploy a postgresql instance to Google Compute
- configure sources against google cloud services
- configure sources against azure cloud services
- configure credentials for all services
- use Vault for secret storage

## Requirements

### Disable On Prem

`onprem.tf` is used internally at Blue Medora. You should disable it with
`mv onprem.tf onprem.tf.old` to avoid trying to deploy it.

### Vault Setup

This deployment depends on vault. Terraform will attempt to read
`VAULT_ADDR` from your environment.

Authenticate against vault with your preferred method:
```
vault login -method=<your method>
or
export VAULT_TOKEN=<your token>
```

### Vault Secrets

Terraform will look for the following secrets:

```
# this is fake data
❯ vault read secret/bindplane/bpcli/poc
Key                     Value
---                     -----
azure_resource_group    mygroup
client_email            lpu-master@project-3d4bda.iam.gserviceaccount.com
client_id               102702109929700103064
client_x509_cert_url    https://www.googleapis.com/robot/v1/metadata/x509/lpu%40myuser-3d4bda.iam.gserviceaccount.com
gcp_project             myproject
postgres_username       user
private_key             rIIEvgIBADANBgkqhkiG9w0BAQEFAASCBKgwggSkAgEAAoIBAQCLLjbry5fCQEid\\nHdsY4Z4J1vV1XWuXwG8Gz/DqOb22J9wK6EvLavBs7lkg1utoU8iVf080hBbeOsDD\\nEFtuwo18aEIdxf3AHV22wFOHPoKnZlJESkKWGvqHA0ITHTgpAtNgHnws0ngKPwTF\\ndRCnCudx4m3AOeh0rnWnphDiUpgGpQAbX5wpGageh9dy5XmQFAo4WyZihsms+yue\\nUh07V/SncRN44JVp0dY7mLKFR4cGVLYmTpUehuWnvr7eI4nJ9rRHExe1tVikpr76\\n88Zj/U9KVO0r8eG12+ZYWcZdT+qzN3AVRvsIDtnh3qZcCde+RfjcW0jF5a/VPrmj\\nHRQ30ErNAgMBAAECggEAMAqIN5TjcdAZoG4FSg3arL/Poy7XbB6m1D2jhV3f74fL\\nqtIrE3B6w8bz6eN1h2HgK0Yx80ki0ZuLHOnA/bbW+pnMNJW6dH1Ocz3otxarJ5go\\njlzppgFy93Z28L0VvQY2KwfqydfuSm8dOQEi+d3IKiLXylHSvK/ZecBXNJ/YzPXW\\nBiRT33htq7CWsEGE0QuUkEBfSSrXx1u03MHrOsKEOnlj5dxfI0TCIWNiB1oPOkln\\nGTi+YIURzNlFjycK2oGGFKnYL4RhlkwaVVS+D+0T2VriMSDbriBT9m1UBTzEfC2S\\nTzBVifnHoUtQZeowMkLp91jsDTGHVktuPSIsLNlHCwKBgQC+fNSj8d8DCWfAOd5c\\ngtLzu3Tr+14/FLpwN3WCdJATw/M2ARPJzj5j8Z+viZcmp+IuaBE3WEsNiT0K1pi6\\nnFgndWYVDuaQgIahC4mZ3M7JqXmUqiJojdc1ZmSWs30g+cPm0YdlMhcpNG7XxtZT\\nYD8uOnPzEZ/7Bjw+0xu8qrqoVwKBgQC7DCKYhjdMlyr8AJVkN0YyMtrnx47VVswj\\nuwa5sn+fpe35oN9PNpwdvBzKTQyVcLA5b5/mQBkQgWErOpFjzhupF3SKuYquvSRB\\n4iOx54ZL7lv1ijOA4sKeKYlJuwbxqVu+N9GpPYosAXU+EP5IPg5CCYxyNDNeBepF\\n46+ZJGw/ewKBgFjz2iik7ktwvO5bF6eDwBbpVvRL8frrJxT4EPvWiuFwA8cYQbFf\\nimsJjlReoCMBCvI4zrFVnda4W7UP+UpLcC8c94ql5q1cF4Jk7ODY6AfsCEaQfHlO\\no0zgf+CP+MLJX94NwnhTJ9WqEojY5YUR0O85hKPhex+yDbgYxT8ZSwkJAoGBAIjp\\nqg0g+Stb33/UYYWYnA40gV11Cg0I2qYwyQx2JsrSJy443hxaac2uGxjNay+b67Iv\\nfcj5FB+rxFdjKHb4r/CGlazRgTzEf9ylzeD0Cq5by/4f6fEmirRAzRgmCUAs6lWD\\nADm0LQZnDs2enLJ+kesumBokMZFaHRCJR8h+C6ovAoGBAKlluMvBljVHyEYMzqzt\\nyqjRU0aUiv5E00Fkri5mKOFre6saLI31SLSOfdi/eRap3wj7ucgmukbqWelXb63d\\nE5+NyX3Mcbb9Aif+77b5/BxehSjidtleBEHHoBua6ZhIXQe3O3Kv1gLiSaBIqevz\\neeSxzoEnSz51xQ04O++Qmkbe
private_key_id          29679bdb50427439h4d731d5db46g1d1ac47d111
project_id              myprojectid
secret_key              e223ccf9-3e93-5f62-ra6q-2a9003b41bc5
```

Notice that the `private key` has all `\n` characters escaped like this : `\\n`

```
❯ vault read secret/azure/account
Key                      Value
---                      -----
client_id                0982e6c3-7e6r-4165-a917-4Tf2cc2c4cd5
client_secret            hZmWPq/IOalnvg/ddrqwTf2Fs342sWQV/sCE90eF0/R=
password                 securepassword
subscription_id          09373b6b-bc8b-4097-925t-eb17334c7d51
tenant_id                75c1e744-e474-4287-946r-de2148dd5ecb
username                 service@company.com
```

## Usage

Make sure you are authenticated to GCP with with a service
account or your personal account
```
gcloud auth
```

Make sure you are authenticated to Azure:
```
az login
```


Set `BINDPLANE_API_KEY`
```
export BINDPLANE_API_KEY=<your api key>
```

Run `terraform`:
```
terraform init
terraform apply
```

Cleanup with:
```
terraform destroy
```
