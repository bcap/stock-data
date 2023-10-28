# first converts all keys to lowercase
recurse(ascii_downcase; .) |

# transform the json array into a series of json objects
.[]