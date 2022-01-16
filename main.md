1. Create a PageFetchJob for the root page
2. Fetch it.
2. Create PageFetchJob with the parent indexing size and the id for the children and put in the queue.
3. Have a repository of nodes 

main
fetch page
are there more pages? fetch them
if there are no more pages to fetch then sum up the indexing sizes and return that value

pagefetchjob {
    id
    parentIndexingSize
}

fetch page
just fetch the page

page {
    size
    indexingSize
    children
}

Also, remember to use sessions or something like that to prevent the connections from being too expensive.