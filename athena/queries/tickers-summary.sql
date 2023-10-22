select 
    "exchangegroup",
    "type",
    count(1) as "amount"
from tickers 
group by
    "exchangegroup",
    "type"
order by
    "exchangegroup" asc,
    "type" asc