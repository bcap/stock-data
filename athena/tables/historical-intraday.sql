DROP TABLE IF EXISTS historical_intraday_1m;

CREATE EXTERNAL TABLE historical_intraday_1m (
    `timestamp` timestamp,
    `ticker` string,
    `low` float,
    `high` float, 
    `open` float,
    `close` float,
    `volume` int
)
ROW FORMAT SERDE 'org.openx.data.jsonserde.JsonSerDe'
LOCATION 's3://bcap-stock-data/historical-intraday/interval-1m/';


DROP TABLE IF EXISTS historical_intraday_1m_plain_text;

CREATE EXTERNAL TABLE historical_intraday_1m_plain_text ( 
	data string
)
ROW FORMAT SERDE 'org.apache.hadoop.hive.serde2.lazy.LazySimpleSerDe'
STORED AS INPUTFORMAT 'org.apache.hadoop.mapred.TextInputFormat' 
OUTPUTFORMAT 'org.apache.hadoop.hive.ql.io.HiveIgnoreKeyTextOutputFormat'
LOCATION 's3://bcap-stock-data/historical-intraday/interval-1m/';