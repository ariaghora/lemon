L = Lemon

L.RGBA = function(r, g, b, a)
    local color = {}
    color.r = r
    color.g = g
    color.b = b
    color.a = a
    return color
end

L.states = {}
L.set_global = function(key, val)
    L.states[key] = val
end

L.get_global = function(key)
    return L.states[key]
end
