select 
    date_format(date_trunc('day', "timestamp"), '%Y-%m-%d - %a') as "day", 
    ticker, 
    sum(volume) as "total_volume",
    min(low) as "min_low",
    max(high) as "max_high"
from 
    historical_intraday_1m
group by 
    date_trunc('day', "timestamp"), 
    ticker
order by 
    "day" asc,
    ticker asc
