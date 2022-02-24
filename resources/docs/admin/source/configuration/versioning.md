# Bucket Versioning

It is generally recommended to use S3 buckets with versioning enabled. This will help protect data from accidental
deletion by authorized users. While this added benefit is very important, configuring bucket lifecycle rules is vital to control
costs due to accumulating versions. 

Note, object version support is not explicitly exposed via the Hoss file browser UI or the `hoss-client` library,
but direct S3 access or other tools can take advantage of versioned data if desired.

Administrators can restore a deleted file by deleting the "Delete Marker". This will restore the file in the
bucket and create an event that the Hoss will detect, resulting in metadata getting re-indexed, and if configured,
the object will be synced.

## Enabling Bucket Versioning
If manually configuring AWS infrastructure, enable bucket versioning in the AWS console.

If using the [https://github.com/WyssCenter/terraform-hoss-aws](https://github.com/WyssCenter/terraform-hoss-aws)
module, set `versioning=true` for the desired buckets. The module ignores the lifecycle rules on the bucket
so you are free to modify those manually as discussed below.

## Configuring Bucket Lifecycle Rules

It is up to an administrator to manually configure lifecycle rules. This can easily be done via the AWS Console.
While recommended rules are outlined below, you are free to configure the rules how ever it best fits your use case.

The recommended behavior is:
* Up to 1 non-current object is kept. Additional versions are removed to minimize storage.
* After a specified time, all non-current objects are removed. This time is the "restore" window.
* If a delete marker has expired it is removed (a delete marker will expire if there are no non-current objects below it)

These features together provide a system where storage of versions is minimized while still allowing for a time period
where objects can be "undeleted" by removing their delete marker, promoting the latest non-current version to the current
version.

This can be practically implemented with two lifecycle rules outlined in the subsections below.

### Reduce Non-current Objects Rule
This rule is added to reduce storage and delete extra non-current objects. If you wish to actually use your versioned objects
you will want to change or skip this rule, as it will only maintain 1 non-current object.

The rule will delete any noncurrent object after the object has been non-current for more than 1 day. It will retain the newest 1 versions, which
means the newest version will not get removed, regardless of how old it is. In the case where you delete an object, the version at the time
of delete becomes this first noncurrent version with a delete marker becoming the current version.

1) In the S3 console, select your bucket
2) Navigate to the "Management" tab
3) Click on "Create lifecycle rule"
4) Fill out the form
   1) Enter a name (e.g. `reduce-noncurrent`)
   2) For the "rule scope", select "Apply to all objects in the bucket".
   3) Under "Lifecycle rule actions" select **only** "Permanently delete noncurrent versions of objects"
   4) Under "Permanently delete noncurrent versions of objects", set "Days after objects become noncurrent" to `1`. 
   5) Under "Permanently delete noncurrent versions of objects", set "Number of newer versions to retain" to `1`. 
   6) Click save



### Permanently Delete Objects Rule
This rule is added to complete a permanent delete of an object after a fixed period time. This is accomplished by removing **all** non-current versions
older than a specified time and deleting expired delete markers. Once all non-current versions are deleted, the delete marker expires and is removed as well.

This rule also cleans up multi-part uploads that have not completed. If a multipart upload is halted in the middle of the process, you are billed for the parts
but cannot access them until the upload is either completed or aborted. Since we expect our uploads to complete in a timely manner, if any parts are around for more
than a day they have been abandoned and should be cleaned up to reduce costs.

1) In the S3 console, select your bucket
2) Navigate to the "Management" tab
3) Click on "Create lifecycle rule"
4) Fill out the form
   1) Enter a name (e.g. `permanent-delete`)
   2) For the "rule scope", select "Apply to all objects in the bucket".
   3) Under "Lifecycle rule actions" select "Permanently delete noncurrent versions of objects" and "Delete expired object delete markers or incomplete multipart uploads"
   4) Under "Permanently delete noncurrent versions of objects", set "Days after objects become noncurrent" to what ever period you wish to support recovery, e.g. `7`.
   5) Under "Permanently delete noncurrent versions of objects", leave "Number of newer versions to retain" blank. 
   6) Under "Delete expired object delete markers or incomplete multipart uploads", check "Delete expired object delete markers"
   7) Under "Delete expired object delete markers or incomplete multipart uploads", check "Delete incomplete multipart uploads" and set "Number of days" to `1`. 
   9) Click save