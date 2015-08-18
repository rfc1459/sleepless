Sleepless
=========

Ugly workaround for an ugly bug.


Seriously, why?!?
-----------------

Because my Toshiba USB3.0 HDD is broken beyond repair and it will go to sleep
after about 20 seconds of inactivity. No amount of `hdparm` or `sdparm` can
fix the thing, so the only alternative is writing something before the timeout
expires.


FAQ
---

Q: Can't you just ask Toshiba to fix it?  
A: According to them, this is a "feature" to save power so there is nothing to
fix. The fact that OS caches might still be dirty and some xHCI controllers
might evict the device altogether from the bus pales in comparison of the few
mW saved! Go green!

Q: What's the problem with the disk spinning down periodically? I mean, come
on!  
A: See the previous question: some xHCI controllers react to the device
powering down by evicting it from the bus. To be honest, the device-side
interface should **NEVER** power down, but apparently Toshiba is really
hell-bent on saving power.

Q: Ok, I get it. But why the whole `O_DIRECT|O_SYNC` mess?  
A: The OS cache can (and will) get in your way.

Q: It doesn't work with NTFS or FAT32!  
A: INVALID/WONTFIX (sounds familiar, doesn't it?).


License
-------

Same, boring old 2-clause BSD. See [LICENSE][] for details.


Contacts
--------

Issues are not welcome (works for me), PRs are welcome _iff_ they make sense
and they don't break my use case.

If your PR gets rejected and you badly need your feature, feel free to fork.


[LICENSE]: https://github.com/rfc1459/sleepless/blob/master/LICENSE
