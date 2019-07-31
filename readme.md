# listpick
Starting as a powershell history selector, this became a generic text selector. It reads from stdin, and produces a termui-based list. Using vim keybindings, navigate up and down the list and visually select the lines you want to keep. Deselect lines using q. Pressing enter will output the lines to stdout, while pressing q will quit without producing output.

Think of it like a homemade dmenu
