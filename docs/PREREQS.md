# Prerequisites

You will need an established Red Hat Partner Connect account. This tool will not
create this for you. 

You'll need to set two environment variables to use `productctl`. 

| Env Var | Description |
|-|-|
|`CONNECT_ORG_ID` |The organization you're working against. Helps in filtering queries.|
|`CONNECT_API_TOKEN`| Your API token. Used to scope requests just to your project.|

### Getting an API token:

Log into Red Hat Partner Connect and access this URL:
https://connect.redhat.com/account/api-keys

### Getting your ORG ID

Log into Red Hat Partner Connect and access this URL:
https://connect.redhat.com/account/company-profile. Your ORG ID should be listed
at the top of this UI.