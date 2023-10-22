# transforms an object into an array with only their values
# Eg:
# {
#   "0": {
#     "Code": "0R2V",
#     "Exchange": "IL",
#     "Name": "Apple Inc."
#   }
# }
# when passed through value_only_array:
# [
#   {
#     "Code": "0R2V",
#     "Exchange": "IL",
#     "Name": "Apple Inc."
#   }
# ]
def value_only_array:
    if . then [to_entries[].value] end
;


# converts a map of objects into an array of the object values with their key flattened back into the object
# Eg:
# {
#   "1994-09-30": {
#     "epsActual": 0.0539
#   }
# }
# when passed through flatten_key("date"):
# [
#   {
#     "date": "1994-09-30",
#     "epsActual": 0.0539
#   }
# ]
def flatten_key($dest_key):
    if . then
        [
            to_entries[] |
            (
                [{key: $dest_key, value: .key}] |
                from_entries
            ) + .value 
        ]
    end
;