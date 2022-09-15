# Basics

> We use Lua programming language to work with Lemon.
If you are not familiar with Lua, we recommend you to check [this crash course by Tyler Neylon](https://tylerneylon.com/a/learn-lua/).
Don't worry, it is a small language and really easy to learn!

## Hello world

Create a project folder with following structure:
```
hello_world/
├── game.json
└── main_scene.lua
```

Prepare project config `game.json`.
Set the starting scene to `main_scene`.
The scene name corresponds to a lua script with the same filename and `.lua` extension.

```json
{
    "title": "Hello world",
    "starting_scene": "main_scene",
    "screen_width": 800,
    "screen_height": 500,
    "target_fps": 60
}
```

Create `main_scene.lua` and add two global functions `on_load()` and `on_update(dt)`.

```lua
function on_load()
end

function on_update(dt)
end
```
The `on_load()` function is where we initialize variables, loading resource.
The `on_update(dt)` is where we put the game logic that will be evaluated frame-by-frame.
These are two minimum functions 

We can run our project by running
```
$ lemon run project_dir/
```
and if everything goes well, we should see this window with a blank white screen:

![hello world screenshot](basics/hello%20world%20screenshot.png)

Congratulations for running your first "game".
Sadly, it looks boring.
Let's add more spices.

## Adding sprites

Sprites are one of main parts in Lemon.
We can use sprites to represent players, enemies, background, etc.
A sprite stores various information such as x and y coordinates, rotation, scale, and texture (a drawable "image" loaded on GPU).
Sprites can be either static or animated.

For a starter, let's show a simple sprite.

```lua
local ball

function on_load()
    ball = L.new_sprite("Ball")
    -- or:
    --   ball = Lemon.new_sprite("My Bunny")
    -- note:
    --   sprite name (e.g., "Ball") must be unique in current scene

    ball:set_texture_from_file("ball.png")
    ball.x, ball.y = 400, 250
    ball.width, ball.height = 64, 64

    -- Setting sprite width and height manually will not automatically update
    -- sprite's origin, so we restore the origin back to center.
    -- We don't need to do this if we don't change sprite's size
    ball:set_origin_center()
end

function on_update(dt)
    -- Move and rotate
    ball.x = ball.x + 400 * dt
    ball.rotation = ball.rotation + 500 * dt

    -- Prevent going beyond the screen
    if ball.x > L.get_screen_width() then
        ball.x = 0
    end
end
```

In Lemon, you will get an access to the global variable `Lemon` (or its alias, `L`).
It provides many functionalities such as sprite creation in above example.

The code is simple and self-explanatory: when the screen is loaded, create a new sprite, set the texture, and set the initial position.
On frame update, also update the x position and do rotation clockwise.
Finally, when it goes beyond screen width, return it's x position back to 0.

> Any sprite will be drawn right away as soon as it is created and its texture is determined.

> Note that we multiply sprite's velocity and rotation amount with `dt`, i.e., the _frame time_.
> We do this to get a visually consistent velocity across different machine with different computing performance.
> Some computers can achieve 60fps, but some others may not.
> Without multiplying with `dt`, we will see much slower movement in slower machines.

When you run the project and if everything is okay, you will see a ball rolling:

![](basics/hello%20world%20screenshot2.png)

## Distributing our game

Currently we only support `"standalone"` mode that includes all scripts and resources into the final executable.
First of all, in `game.json` set `"build_mode"` to `"standalone"`.

```json
{
    "title": "The Basics",
    "starting_scene": "main_scene",
    "screen_width": 800,
    "screen_height": 500,
    "target_fps": 60,
    "build_mode": "standalone"
}
```
Then execute following command:

```bash
$ lemon build project_dir/
```

You will find the executable file under the project directory with the same file name as the game title.