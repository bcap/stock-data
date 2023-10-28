# first converts all keys to lowercase
recurse(ascii_downcase; .) |

# transform the json array into a series of json objects
.[] |

# operatingmic can be a string encoded array of items or even null
.operatingmic |= if . then split(", ") end