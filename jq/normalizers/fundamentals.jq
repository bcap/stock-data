# The fundamentals api responses are a bit weird in term of data structure
#
# This normalizer will:
#
# 1. Convert json objects that represent arrays into actual arrays. For instance:
# {
#   "0": { "a": 1 },
#   "1": { "b": 2 }
# }
# will be transformed with the `value_only_array` function into:
# [
#   { "a": 1 },
#   { "b": 2 },
# ]
#
# 2. Converts a map of objects into an array of the object values with their key flattened back into the object. For instance:
# {
#   "1994-09-30": {
#     "epsActual": 0.0539
#   },
# }
# will be transformed with the `flatten_key` function into:
# [
#   {
#     "date": "1994-09-30",
#     "epsActual": 0.0539
#   }
# ]


# first converts all keys to lowercase
recurse(ascii_downcase; .) |


.general?.listings?                         |= value_only_array |
.general?.officers?                         |= value_only_array |

if .earnings? then
    .earnings?.annual?                      |= flatten_key("date") |
    .earnings?.history?                     |= flatten_key("date") |
    .earnings?.trend?                       |= flatten_key("date")
end                                         |

if .esgscores? then
    .esgscores.activitiesinvolvement?       |= value_only_array
end                                         |

if .etf_data? then
    .etf_data.asset_allocation?             |= flatten_key("type") |
    .etf_data.world_regions?                |= flatten_key("region") |
    .etf_data.sector_weights?               |= flatten_key("sector") |
    .etf_data.fixed_income?                 |= flatten_key("type") |

    # we dont need top 10 holdings when we have all holdings
    .etf_data.holdings?                     |= value_only_array |
    del(.etf_data.top_10_holdings)          |

    # valuations_growth seems to try to put 2 different types of keys (valuations and growth) into the same valuations_growth key
    # here we split those 2 back into 2 different keys
    .etf_data.valuations = [
        {type: "portfolio"} + .etf_data?.valuations_growth?.valuations_rates_portfolio?,
        {type: "category"} +  .etf_data?.valuations_growth?.valuations_rates_to_category?
    ] |
    .etf_data.growth = [
        {type: "portfolio"} + .etf_data?.valuations_growth?.growth_rates_portfolio?,
        {type: "category"} +  .etf_data?.valuations_growth?.growth_rates_to_category?
    ] |
    del(.etf_data.valuations_growth)
end                                         |

if .financials? then
    .financials?.balance_sheet?.quarterly?      |= flatten_key("date") |
    .financials?.balance_sheet?.yearly?         |= flatten_key("date") |
    .financials?.cash_flow?.quarterly?          |= flatten_key("date") |
    .financials?.cash_flow?.yearly?             |= flatten_key("date") |
    .financials?.income_statement?.quarterly?   |= flatten_key("date") |
    .financials?.income_statement?.yearly?      |= flatten_key("date")
end                                             |

# merge .holders.funds and .holders.instituitions into a single .holders array
if .holders? then
    .holders?                               |= (
        if .funds?        then [ .funds        | value_only_array[] | {type: "fund"}        + .] else [] end +
        if .institutions? then [ .institutions | value_only_array[] | {type: "institution"} + .] else [] end
    )
end                                         |

if .insidertransactions? then
    .insidertransactions                    |= value_only_array
end                                         |

if .outstandingshares? then
    .outstandingshares.annual?              |= value_only_array |
    .outstandingshares.quarterly?           |= value_only_array
end                                         |

if .splitsdividends? then
    .splitsdividends?.numberdividendsbyyear?    |= value_only_array
end                                             |

.