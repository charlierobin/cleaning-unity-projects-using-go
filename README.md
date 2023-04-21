# A suite of compiled Go command line apps for detecting and cleaning up Unity project directories
 
A Go version of the Python Unity projects cleaner script ([here](https://github.com/charlierobin/clean-unity-projects)).

For a bit of additional interest, I split the app up into constituent apps, each of which carries out a particular function. (As opposed to keeping it all in one big monolithic thing, which is what the Python script is.)

The main app is [unity-projects-cleaner](https://github.com/charlierobin/cleaning-unity-projects-using-go/tree/main/unity-projects-cleaner)

This presents a list of places/mounted volumes on your computer, allows you to select one or more of them, then scans those locations, detects Unity projects folders, lists them, then gives you the option of “cleaning” them, ie: moving the `Logs`, `Library` and `obj` folders to the trash.

On my main work drive, with 85 Unity projects, most of which I hardly look at any more, this freed up **a lot** of extra free space.

So, the up side: lots of new extra free space.

The down side: if you ever come back to open any of those projects, it will take a bit longer because Unity will have to regenerate a load of cached stuff. On a big project, that can take quite a bit of time. *What the heck, for me it was worth it.*

https://user-images.githubusercontent.com/10506323/233540403-978f67d2-2a73-4555-b9b9-2c4f0ce805e5.mp4

The **main** components are:

## unity-projects-cleaner

app calls the other apps via Go’s `os/exec` package:

```
cmd := exec.Command("shellCommandHere", arg1, arg2, arg3)
stdout, err := cmd.Output()
```

etc.

They all return their results via `stdout`.

The way the functionality has been split up is all a bit contrived, but apart from having another look at Go, it was just something I was interested in experimenting with.

The components are:

### list-mounted-volumes

Does exactly what it says: returns a list of the user’s mounted volumes.

Under the hood, it just wraps `df` and filters out most of what is returned (the macOS Recovery volume and a load of other stuff) because all I was really interested in were any mounted external volumes, not all the internal system rubbish (which on my system was "/", "/dev", "/System/Volumes/Data", "/private/var/vm", "/System/Volumes/Data/home", "/Volumes/Recovery").

The app adds the user’s home documents directory to the list of external volumes derived from above, and allows the user to selected one of more places to be scanned. (If nothing is selected, the app exits.)

Each selected place is in turn passed to …

### find-unity-project-folders

This uses `filepath.Walk` to look for Unity project folders.

Each candidate path is passed to …

### get-unity-project

… which examines it to see if it’s a Unity project.

A Unity project folder is defined as a folder containing `Assets`, `Packages`, `ProjectSettings` and `UserSettings` subdirectories, with a `ProjectVersion.txt` file and a `ProjectSettings.asset` file in the `ProjectSettings` folder.

If found, some data is extracted and returned: the Unity editor version of the project, the product name and the company name. (As set in the Player section of the Project Settings editor window.)

Once we have a list of project folders, they can be reviewed just to check that all seems well.

If the user then elects to go ahead with a clean, each project folder path is sent to:

### clean-unity-project

This actually deletes the three cache folders (`Logs`, `Library` and `obj`), using:

### trash

This also uses `exec.Command`, calling `osascript` to instruct **the Finder** to move the specified item to the user‘s trash.

It appends a note to the folder name so it can still be identified once in the trash, just in case something goes wrong.

So the Logs folder in a project called “Asteroids” will end up in the trash named “Logs (from Asteroids)”.

I used this method instead of just using `rm` (or just directly deleting the file using Go’s filesystem commands) partly to see how else it could be done, but also because I wanted a last line of defence just in case of error. Doing it this way means that there’s a chance to review the deleted folders in the trash before emptying it, and also the Finder’s “Put Back” command can be used to replace any item back in its original place with ease.

The drawback of this is that the `trash` app will **possibly** need to be granted accessibility permissions, depending on what version of macOS / OSX you are using.

### Finally …

Because all the deleted folders are not actually deleted but are only moved to the trash, you won’t actually get hard drive space back until you commit yourself to emptying your trash.

## Building

In Terminal, cd to the project directory, then run [build.command](https://github.com/charlierobin/cleaning-unity-projects-using-go/blob/main/build.command)

This compiles all the components for both arm64 (aka Apple Silicon, the new M chips) and amd64 (ie: Intel 64 bit).

The command then runs `mv` to move the amd64 binaries up to the **top level of the build directory**, because that’s where I’ve been running them from, and that’s where they expect to be. (Hence the presence of `config.json` in this directory, which we'll get to later.) You will need to edit this if your Mac uses Apple’s new chips.

## Variations on a theme

If all that seems a bit too involved, and you prefer doing more stuff by yourself, you can call `find-unity-project-folders` directly from the command line, just pass in the path where you want the find to start.

The first parameter to `find-unity-project-folders` is always the path to be checked. After that, passing in -v turns on the verbose logging option (to a log file created in the build directory). Any other paths passed in are treated as locations to skip, which helps in speeding the scan process up. (I have a large `Applications` folder on the same drive I was checking, as well as a large fonts folder, so this came in handy.)

Alternatively, you can add paths to be skipped to the `config.json` file, also in the builds folders.

Once you have your list of Unity projects, you can call `clean-unity-project` on a case by case basis:

`clean-unity-project /path/to/unity/project/here`

`clean-unity-project` calls `get-unity-project` just to check that the path you have passed is indeed a Unity project.

You can call `get-unity-project` yourself to check whether or not the app will identify the folder as a Unity project:

`get-unity-project /path/to/unity/project/here`

If it thinks it’s a Unity project, you’ll see the path, the Unity editor version, product name, and company name output.

Otherwise you’ll see nothing.

## Notes

With the exception of the implementation of `trash` and `list-mounted-volumes`, I think most of this would be happy enough on Windows with a bit of tweaking here and there. (See the additional added notes below.) But having said that, I’ve not actually tested any of it out.

---

### Additional observations on the above…

The above is perhaps not really true anymore, as I’ve added `picker`, which is a generic app for selecting from a bunch of options. It’s used both to allow the user to select from the list of available volumes, and to select from the list of found projects.

It’s opened by Applescript commands sent to the Terminal app, because I wanted it to appear in a new tab, whilst the calling app waits in the background. (There was no particular compelling design decision behind this, I was just curious to see how it could be made to work.)

Which means that to get it working on Windows or Linux you’ll have to deal with the above, as well as the couple of other places where it’s all gone a bit macOS-centric.

---

At the moment, everything is designed to be run from the top level of the `build` directory. (That’s where the `config.json` and log file live, and that’s where all the apps will expect to find each other when they try to call each other. There’s no adding these apps to your bash path or anything like that.)

There are other alternatives to `trash`, like this one:

https://github.com/ali-rantakari/trash

… or …

https://github.com/morgant/tools-osx

… which could quite easily be substituted in if you preferred.

`stdout` is used for everything, all errors are output to the log file.


