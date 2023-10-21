# transform the json array into a series of json objects
.[] |

# OperatingMic can be a string encoded array of items or even null
.OperatingMIC |= if . then split(", ") end