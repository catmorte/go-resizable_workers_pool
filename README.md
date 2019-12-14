# go-resizable_workers_pool

To create new pool:

    initialPoolSize := uint(1)
    wp := wpool.NewPool(context.Background(), initialPoolSize)


or to be able to kill the entire pool:

    initialPoolSize := uint(1)
    ctx, cancel := context.WithCancel(context.Background())
    wp := wpool.NewPool(ctx, initialPoolSize)
    ...
    cancel()
    

to add new job to pool:

    wp.Do(func(){
        ... 
    })
    
to resize pool:

    initialPoolSize := uint(1)
    wp := wpool.NewPool(context.Background(), initialPoolSize)
    ...
    err = wp.Resize(100)
    ...
    err = wp.Resize(1)
    
    