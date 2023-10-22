DROP TABLE IF EXISTS tickers;

CREATE EXTERNAL TABLE tickers (
    `Code` string,
    `Country` string,
    `Currency` string,
    `Exchange` string,
    `Exchange` string,
    `Isin` string,
    `Name` string,
    `Type` string
)
ROW FORMAT SERDE 'org.openx.data.jsonserde.JsonSerDe'
LOCATION 's3://bcap-stock-data/tickers/';


DROP TABLE IF EXISTS tickers_plain_text;

CREATE EXTERNAL TABLE tickers_plain_text (
	data string
)
ROW FORMAT SERDE 'org.apache.hadoop.hive.serde2.lazy.LazySimpleSerDe'
STORED AS INPUTFORMAT 'org.apache.hadoop.mapred.TextInputFormat'
OUTPUTFORMAT 'org.apache.hadoop.hive.ql.io.HiveIgnoreKeyTextOutputFormat'
LOCATION 's3://bcap-stock-data/tickers/';