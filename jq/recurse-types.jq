
def recurse(level):
    # debug((level | tostring) + ": " + (. | tostring)) |
    if type == "object" then 
        [ to_entries[] | {key: .key, value: .value | recurse(level+1)} ] | from_entries
    elif type == "array" then
        if (. | length) == 0 then []
        else [.[0] | recurse(level+1)]
        end
    else
        type
    end
;

recurse(0)