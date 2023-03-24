# Known issues to work on or fix

1. the notion task monitor can pick up to-do while the user is still writing on the page, thus creates a incomplete Task, (and later another complete Task). Should we only scan page when it is "closed" (not sure whether notion API support)

2. the notion taks monitor can pick up the to-do, and create a Task, but the title does not tell any context, for example, in a page/task "1-on-1 with Joe", user may have a to-do "follow up of the deployment improvements". Ideally, the next Task title will be something like "Joe - follow up of the deployment improvements". So, a little bit of intelligence would be nice here.