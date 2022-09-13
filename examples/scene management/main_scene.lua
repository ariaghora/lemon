function on_update(dt)
    L.draw_rect_fill(
        0,
        0,
        L.get_screen_width(),
        L.get_screen_height(),
        L.rgba(255, 123, 0, 255)
    )

    L.draw_text(
        "This is main scene\n" ..
        "Press backspace to go back to the title scene",
        0,
        0,
        30,
        L.rgba(255, 255, 255, 255)
    )

    if L.is_key_pressed(L.KEY_BACKSPACE) then
        L.set_scene("title_scene")
    end
end
