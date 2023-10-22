select 
    "country",
    "type",
    count(1) as "amount"
from tickers 
group by
    "country",
    "type"
order by
    "country" asc,
    "type" asc