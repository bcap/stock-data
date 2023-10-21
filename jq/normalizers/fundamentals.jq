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

.Earnings?.Annual?                          |= flatten_key("date") |
.Earnings?.History?                         |= flatten_key("date") |
.Earnings?.Trend?                           |= flatten_key("date") |

.ESGScores?.ActivitiesInvolvement?          |= value_only_array |

.Financials?.Balance_Sheet?.quarterly?      |= flatten_key("date") |
.Financials?.Balance_Sheet?.yearly?         |= flatten_key("date") |
.Financials?.Cash_Flow?.quarterly?          |= flatten_key("date") |
.Financials?.Cash_Flow?.yearly?             |= flatten_key("date") |
.Financials?.Income_Statement?.quarterly?   |= flatten_key("date") |
.Financials?.Income_Statement?.yearly?      |= flatten_key("date") |

.General?.Listings?                         |= value_only_array |
.General?.Officers?                         |= value_only_array |

.Holders?.Funds?                            |= value_only_array |
.Holders?.Institutions?                     |= value_only_array |

.InsiderTransactions?                       |= value_only_array |

.outstandingShares?.annual?                 |= value_only_array |
.outstandingShares?.quarterly?              |= value_only_array |

.SplitsDividends?.NumberDividendsByYear?    |= value_only_array |

.