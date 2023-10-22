select
    general.code,
    holder
from fundamentals
CROSS JOIN UNNEST(holders) as t(holder)
where general.code in ('MSFT', 'AAPL')
order by
    general.code,
    holder.currentshares desc