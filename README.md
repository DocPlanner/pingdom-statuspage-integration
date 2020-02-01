# pingdom-statuspage-integration
## How it works?
pingdom-statuspage-integration finds StatusPage Component based on name of Pingdom Check and changes it's status based on `currentStatus` field from Pingdom Webhook. 
If there is more than one StatusPage Component with the same name(f.ex. on multiple pages) as Pingdom Check then status of all those components will be affected.

## Adding new components on StatusPage
State of StatusPage components is kept inside of application and refreshed every 30 minutes. It means that when you add Pingdom Check and corresponding to it StatusPage Component you need to restart application or wait up to 30 minutes to state refresh. 

## Pingdom configuration
Only thing to do is to add webhook to your check in Pingdom
### Sample webhook url
`https://your.domain.tld/?secret=SECRET`
SECRET is a value defined in environment variable "SECRET"

## Environment variables
### Required
SECRET - secret used in communication from Pingdom
STATUSPAGE_TOKEN - StatusPage API key
### Optional
MAX_RETRIES - number of retries (default: 2)
RETRY_INTERVAL - numer of seconds between retries (default: 10)