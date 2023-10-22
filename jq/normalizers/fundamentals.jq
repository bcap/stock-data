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


.General?.Listings?                         |= value_only_array |
.General?.Officers?                         |= value_only_array |

if .Earnings? then
    .Earnings?.Annual?                      |= flatten_key("date") |
    .Earnings?.History?                     |= flatten_key("date") |
    .Earnings?.Trend?                       |= flatten_key("date")
end                                         |

if .ESGScores? then
    .ESGScores.ActivitiesInvolvement?       |= value_only_array
end                                         |

if .ETF_Data? then
    .ETF_Data.Asset_Allocation?             |= flatten_key("type") |
    .ETF_Data.World_Regions?                |= flatten_key("region") |
    .ETF_Data.Sector_Weights?               |= flatten_key("sector") |
    .ETF_Data.Fixed_Income?                 |= flatten_key("type") |

    # we dont need top 10 holdings when we have all holdings
    .ETF_Data.Holdings?                     |= value_only_array |
    del(.ETF_Data.Top_10_Holdings)          |

    # Valuations_Growth seems to try to put 2 different types of keys (Valuations and Growth) into the same Valuations_Growth key
    # Here we split those 2 back into 2 different keys
    .ETF_Data.Valuations = [
        {type: "Portfolio"} + .ETF_Data?.Valuations_Growth?.Valuations_Rates_Portfolio?,
        {type: "Category"} +  .ETF_Data?.Valuations_Growth?.Valuations_Rates_To_Category?
    ] |
    .ETF_Data.Growth = [
        {type: "Portfolio"} + .ETF_Data?.Valuations_Growth?.Growth_Rates_Portfolio?,
        {type: "Category"} +  .ETF_Data?.Valuations_Growth?.Growth_Rates_To_Category?
    ] |
    del(.ETF_Data.Valuations_Growth)
end                                         |

if .Financials? then
    .Financials?.Balance_Sheet?.quarterly?      |= flatten_key("date") |
    .Financials?.Balance_Sheet?.yearly?         |= flatten_key("date") |
    .Financials?.Cash_Flow?.quarterly?          |= flatten_key("date") |
    .Financials?.Cash_Flow?.yearly?             |= flatten_key("date") |
    .Financials?.Income_Statement?.quarterly?   |= flatten_key("date") |
    .Financials?.Income_Statement?.yearly?      |= flatten_key("date")
end                                             |

# Merge .Holders.Funds and .Holders.Instituitions into a single .Holders array
if .Holders? then
    .Holders?                               |= (
        if .Funds?        then [ .Funds        | value_only_array[] | {type: "fund"}        + .] else [] end +
        if .Institutions? then [ .Institutions | value_only_array[] | {type: "institution"} + .] else [] end
    )
end                                         |

if .InsiderTransactions? then
    .InsiderTransactions                    |= value_only_array
end                                         |

if .outstandingShares? then
    .outstandingShares.annual?              |= value_only_array |
    .outstandingShares.quarterly?           |= value_only_array
end                                         |

if .SplitsDividends? then
    .SplitsDividends?.NumberDividendsByYear?    |= value_only_array
end                                             |

.