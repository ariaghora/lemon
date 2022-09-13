local x_acceleration, y_acceleration = 70, 30
local x_deceleration, y_deceleration = 30, 30
local x_velocity, y_velocity = 0, 0
local max_velocity = 600
local bullets_per_second = 3

-- 0 for not moving state, 1 for moving towards positive axis, and
-- -1 for moving towards negative axis
local moving_x, moving_y = 0, 0

function on_load()
    this:set_texture_from_file("assets/images/ship1.png")
    this.x, this.y = 80, 200
end

local function handle_movement(dt)
    if L.is_key_down(L.KEY_RIGHT) then
        x_velocity = x_velocity + x_acceleration * dt
        moving_x = 1
    elseif L.is_key_down(L.KEY_LEFT) then
        x_velocity = x_velocity - x_acceleration * dt
        moving_x = -1
    else
        moving_x = 0
    end

    if L.is_key_down(L.KEY_DOWN) then
        y_velocity = y_velocity + y_acceleration * dt
        moving_y = 1
    elseif L.is_key_down(L.KEY_UP) then
        y_velocity = y_velocity - y_acceleration * dt
        moving_y = -1
    else
        moving_y = 0
    end

    -- cap the velocity to max_velocity
    if x_velocity > max_velocity * dt then
        x_velocity = max_velocity * dt
    elseif x_velocity < -max_velocity * dt then
        x_velocity = -max_velocity * dt
    end

    -- handle deceleration
    if moving_x == 0 and x_velocity > 0 then
        x_velocity = x_velocity - x_deceleration * dt
    elseif moving_x == 0 and x_velocity < 0 then
        x_velocity = x_velocity + x_deceleration * dt
    end

    if moving_y == 0 and y_velocity > 0 then
        y_velocity = y_velocity - y_deceleration * dt
    elseif moving_y == 0 and y_velocity < 0 then
        y_velocity = y_velocity + y_deceleration * dt
    end

    -- cap velocity close to zero into zero to eliminate little movements when
    -- no arrow key pressed
    if moving_x == 0 then
        if x_velocity < x_deceleration * dt and x_velocity > -x_deceleration * dt then
            x_velocity = 0
        end
    end
    if moving_y == 0 then
        if y_velocity < y_deceleration * dt and y_velocity > -y_deceleration * dt then
            y_velocity = 0
        end
    end

    -- update player position according to the updated x & y coords.
    this.x = this.x + x_velocity
    this.y = this.y + y_velocity

    -- prevent player to move out of screen boundary
    if this.x > L.get_screen_width() - this.width / 2 then
        this.x = L.get_screen_width() - this.width / 2
        x_velocity = 0
    elseif this.x < this.width / 2 then
        x_velocity = 0
        this.x = this.width / 2
    end

    if this.y > L.get_screen_height() - this.height / 2 then
        this.y = L.get_screen_height() - this.height / 2
        y_velocity = 0
    elseif this.y < this.height / 2 then
        y_velocity = 0
        this.y = this.height / 2
    end
end

local bullet_num = 0
local counter = 0
local shooting = false
local function handle_firing(dt)
    if L.is_key_down(L.KEY_SPACE) then
        shooting = true
    elseif L.is_key_up(L.KEY_SPACE) then
        shooting = false
        counter = 0
    end

    if shooting then
        if (counter % math.floor(L.get_fps() / bullets_per_second) == 0) or
            (counter == 0) then
            bullet_num = bullet_num + 1
            local b = L.new_sprite("bullet" .. tostring(bullet_num))
            b.x = this.x + 32
            b.y = this.y
            b:set_script("sprites/player_bullet.lua")
        end
        counter = counter + 1
    end
end

function on_update(dt)
    handle_movement(dt)
    handle_firing(dt)
end
