# first converts all keys to lowercase
recurse(ascii_downcase; .) |

# transform the json array into a series of json objects
.[] |

# delete some columns we are not interested or are redundant
del(.datetime) |
del(.gmtoffset) |

.