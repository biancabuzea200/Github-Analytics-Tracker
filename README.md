## How to run locally 

1. Create a .env with the variables below

STAGE=DEV

POSTGRES_DSN=postgres://[you]@localhost:5432/[you]

DB_SCHEMA=[SCHEMA_NAME]


DEV_REPO_OWNER=[you]

DEV_REPO=[YOUR_REPO_NAME]

DEV_REPO_API_KEY=[YOUR_GITHUB_PERSONAL_ACCESS_TOKEN]

2. run: go run main.go  

**Description**

This project aims to collect data on repo usage and traffic from the GitHub api and store it for posterity. This is needed because GitHub only 
saves this data for up to two weeks. This data can be used to measure the success of devrel programs - having data to clearly measure success is very useful. 

The project runs efficiently and can scale to run in near constant time for any number of repos that need to be analyzed. The data is stored on postgres and retrieval can be further enchanced by installing TimescaleDB on top, which allows much more performat retrieval of time series data (which this project uses). The app itself is designed to run a k8s job once per day, week or fortnight. There is a docker file provided which builds a simple image based on the official golang docker image that runs the binary. There is more information provided on intended deployment patterns in the dockerfile as a comment. 


**What does it measure?**

The project currently supports collecting data on; 

<ol>
    <li> clones: how many times your repo was cloned per day </li>
    <li> paths: which parts of your repo have the most activity (views) </li>
    <li> sources: how people are navigating to your repos page (through Google, blogs, marketing campaigns etc.)  </li>
</ol>










