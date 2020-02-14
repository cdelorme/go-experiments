
# v1

The primary goal for this iteration was a streamlined first pass at a single file skeleton.

_To run this code you will need to add an image to `assets/sprites/me.jpg`, or change the file path to a file of your choice as I did not want to commit binary files to this repository._


## features

This particular iteration has the following:

- GLFW with semi-smart window management functionality
- OpenGL initialization
- Vertex Arrays
- Vertex Buffers
- Index Buffers
- Drawing Elements using VA/VB/IB
- Textures

I tend to favor really well designed and convenient or simple window management.

The ability to close with the `escape` key.

The ability to toggle fullscreen with `alt+enter` (relatively standard).

Enable and disable vsync with the `v` key.

The ability to print a rolling average FPS on demand using the `d` key.

Window resize support that retains the aspect ratio; granted this is currently leveraging window hints that require the window manager to acknowledge and respect the settings (there are also curious bugs like OpenBox having off-by-one dimensions).

A potentially wasted computation to retain the aspect ratio by adding black bars by controlling the viewport (_more on this later..._).

A really rough first iteration of using a projection matrix to support aspect ratio and converting from a fixed pixel coordinate system to OpenGL which deals with -1 to 1 floating points.


### concerns

I am torn on how to deal with window resizing.

Depending on what you are building you may want to enforce not only the aspect ratio of what is drawn but also the maximum size of the viewport, and to retain both it just seemed like black bars made sense.

For example, if you are building a competitive game then allowing players with larger displays to have a wider view is unfair to the player that cannot afford the same hardware.  _I am actually not certain how games balance this currently, although in FPS I believe there are limits to how much the brain can process making 4:3 preferred anyways and that would kind of eliminate the concerns here.  Also a strategy game might solve this with a fog-of-war effect and the camera view is unimportant._

I also just recently learned about projection matrices, and I am pretty confident at this point that either the camera or the projection matrices can be used to create a bounding box effect.  _It may also be perfectly fine to retain the aspect ratio using the projection matrix and let the camera size according to the window dimensions rather than scaling, which would make the "scale" a different setting that could be applied via the project matrix._

I am also not sure about the vsync setting (eg. `glfw.SwapInterval`).  A few years ago I recall trying to setup glfw and running into major FPS performance problems caused by this setting.  Supposedly it was that combined with some other setting in Xorg.  However now it seems to work without any issues.  Having it on by default means the FPS computation will always be around 60 making it harder to test performance.  I currently use the `TearFree` option with an AMD GPU and that does not appear to be causing any such FPS restrictions.  This is why I created a toggle, so that I could deal with this on a case by case basis.


## future

I want to test passing an array of structs instead of float32 arrays to OpenGL.

I want to start writing abstractions on OpenGL components to streamline an initialization and rendering loop.

I want to add a scene with parallax background support using projection matrices (eg. instead of multiple translations).

I want to add imgui or nuklear to render an interface.

I want to add controller support to navigate the menu system.

I want to add a configuration menu and save configuration in a simple way that can also be modified raw (eg. json).

I want to investigate adding audio (_golang lacks audio packages_).

There is more, but these seem like a sizable start.  The next iteration will probably be some combination of the items in this list, but probably not all of them at once.


# notes

There are some odd choices in my code as well.

I have it suggesting OpenGL 4.1 as a window hint on initialization.

However, I am using the 4.3 OpenGL function for capturing and printing errors.

The latest version of OpenGL is 4.6, but most cards only support 4.5, and Vulkan seems to be the modern option.

Finally, OSX only has 4.1 support and has said it is deprecated meaning they may not support it at all in the future.

I honestly hate that OSX has chosen to forgo a cross platform option in favor of its incredibly vendor-lockin `Metal` implementation.  If they at least mentioned vulkan support I would feel a bit better about this.
