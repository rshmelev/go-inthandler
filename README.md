go-inthandler
=============

Handling interrupts like Ctrl+C in Golang 

First, import it:  
`import interrupts "github.com/rshmelev/go-inthandler"`

Second, make it monitor interrupts during your app initialization:  
`interrupts.TakeCareOfInterrupts(false)`

Next, you can do `<-interrupts.StopChannel`   
or monitor value of `*interrupts.StopPointer`  
to be sure you're not missing forced shutdown event

You can also start shutdown process by calling:    
`interrupts.InterruptTheApp()`

Modify `interrupts.MaxTimeToWaitForCleanup` to ensure your app has enough time to clean up.  
Module is calling `os.Exit(0)` if your app is running too long after interruption.

Tested on Windows and Linux, should work everywhere.

Thanks!