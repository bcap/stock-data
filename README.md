# stock-data

This is a golang program used to build historical data around stocks / funds / etfs

It uses:
- [EODHD](https://eodhd.com/) as source of truth for financial data
- [AWS S3](https://aws.amazon.com/s3/) for final storage of processed files
- [AWS Athena](https://aws.amazon.com/athena/) for querying the generated data

Today the program supports:
- listing all avaialble exchanges
- listing all tickers for given exchanges
- fetching historical intraday data for particular tickers
- fetching fundamentals about particular tickers

Check [sample-config.yaml](sample-config.yaml) for a bare-bones configuration example

## How

Overall the program works by:
1. reading the passed config file
2. fetching relevant data from eodhd.com
3. transforming and normalizing such data with internal jq scripts (those are executed used a golang jq implementation, so no need for the `jq` tool installed) + golang code
4. upload such data to AWS S3

Check the [athena/](athena/) directory for example queries and table definitions