
## 调用顺序

```C
crtools.c 
    int cr_pre_dump_tasks(pid_t pid)
        pre_dump_one_task()
        ret = cr_dump_shmem()
        irmap_predump_prep()
        cr_pre_dump_finish()
    cr_pre_dump_finish()
        pr_info("Pre-dumping tasks' memory\n");
	    for_each_pstree_item(item) 
            struct page_pipe *mem_pp;
		    struct page_xfer xfer;

            ret = open_page_xfer(&xfer, CR_FD_PAGEMAP, vpid(item));
            if (ret < 0)
                goto err;

            mem_pp = dmpi(item)->mem_pp;

            if (opts.pre_dump_mode == PRE_DUMP_READ) {
                timing_stop(TIME_MEMWRITE);
                pr_info("\t\\__ in cr-dump.c:cr_pre_dump_finish ==> \n\t\t\\__ready to call page-xfer.c:page_xfer_predump_pages \n");

                ret = page_xfer_predump_pages(item->pid->real, &xfer, mem_pp);
            } else {
                pr_info("\t\\__ in cr-dump.c:cr_pre_dump_finish ==> ready to call page-xfer.c:page_xfer_dump_pages \n");

                ret = page_xfer_dump_pages(&xfer, mem_pp);
            }

            xfer.close(&xfer);

            if (ret)
                goto err;
```

获取pagemap

```C
// vma.h 
struct vm_area_list {
	struct list_head h;   /* list of VMAs */
	unsigned nr;	      /* nr of all VMAs in the list */
	unsigned int nr_aios; /* nr of AIOs VMAs in the list */
	union {
		unsigned long nr_priv_pages; /* dmp: nr of pages in private VMAs */
		unsigned long rst_priv_size; /* rst: size of private VMAs */
	};
	unsigned long nr_priv_pages_longest;   /* nr of pages in longest private VMA */
	unsigned long nr_shared_pages_longest; /* nr of pages in longest shared VMA */
};

// cr_dump.c 119
int collect_mappings(pid_t pid, struct vm_area_list *vma_area_list, dump_filemap_t dump_file)
{
	int ret = -1;

	pr_info("\n");
	pr_info("Collecting mappings (pid: %d)\n", pid);
	pr_info("----------------------------------------\n");

	ret = parse_smaps(pid, vma_area_list, dump_file);
	if (ret < 0)
		goto err;

	pr_info("Collected, longest area occupies %lu pages\n", vma_area_list->nr_priv_pages_longest);
	pr_info_vma_list(&vma_area_list->h);

	pr_info("----------------------------------------\n");
err:
	return ret;
}

```


怎么获取pages文件夹

```C
int open_page_xfer(struct page_xfer *xfer, int fd_type, unsigned long img_id)
{
	xfer->offset = 0;
	xfer->transfer_lazy = true;

	if (opts.use_page_server)
		return open_page_server_xfer(xfer, fd_type, img_id);
	else
		return open_page_local_xfer(xfer, fd_type, img_id);
}

```

使用本地文件

```C
static int open_page_local_xfer(struct page_xfer *xfer, int fd_type, unsigned long img_id)
{
	u32 pages_id;

	xfer->pmi = open_image(fd_type, O_DUMP, img_id);
	if (!xfer->pmi)
		return -1;

	xfer->pi = open_pages_image(O_DUMP, xfer->pmi, &pages_id);
	if (!xfer->pi)
		goto err_pmi;

	/*
	 * Open page-read for parent images (if it exists). It will
	 * be used for two things:
	 * 1) when writing a page, those from parent will be dedup-ed
	 * 2) when writing a hole, the respective place would be checked
	 *    to exist in parent (either pagemap or hole)
	 */
	xfer->parent = NULL;
	if (fd_type == CR_FD_PAGEMAP || fd_type == CR_FD_SHMEM_PAGEMAP) {
		int ret;
		int pfd;
		int pr_flags = (fd_type == CR_FD_PAGEMAP) ? PR_TASK : PR_SHMEM;

		/* Image streaming lacks support for incremental images */
		if (opts.stream)
			goto out;

		if (open_parent(get_service_fd(IMG_FD_OFF), &pfd))
			goto err_pi;
		if (pfd < 0)
			goto out;

		xfer->parent = xmalloc(sizeof(*xfer->parent));
		if (!xfer->parent) {
			close(pfd);
			goto err_pi;
		}

		ret = open_page_read_at(pfd, img_id, xfer->parent, pr_flags);
		if (ret <= 0) {
			pr_perror("No parent image found, though parent directory is set");
			xfree(xfer->parent);
			xfer->parent = NULL;
			close(pfd);
			goto out;
		}
		close(pfd);
	}

out:
	xfer->write_pagemap = write_pagemap_loc;
	xfer->write_pages = write_pages_loc;
	xfer->close = close_page_xfer;
	return 0;

err_pi:
	close_image(xfer->pi);
err_pmi:
	close_image(xfer->pmi);
	return -1;
}
```

怎么写pages文件
```C

/* local xfer */
static int write_pages_loc(struct page_xfer *xfer, int p, unsigned long len)
{
	ssize_t ret;
	ssize_t curr = 0;

	while (1) {
		ret = splice(p, NULL, img_raw_fd(xfer->pi), NULL, len - curr, SPLICE_F_MOVE);
		if (ret == -1) {
			pr_perror("Unable to spice data");
			return -1;
		}
		if (ret == 0) {
			pr_err("A pipe was closed unexpectedly\n");
			return -1;
		}
		curr += ret;
		if (curr == len)
			break;
	}

	return 0;
}


```


```C
struct page_xfer {
	/* transfers one vaddr:len entry */
	int (*write_pagemap)(struct page_xfer *self, struct iovec *iov, u32 flags);
	/* transfers pages related to previous pagemap */
	int (*write_pages)(struct page_xfer *self, int pipe, unsigned long len);
	void (*close)(struct page_xfer *self);

	/*
	 * In case we need to dump pagemaps not as-is, but
	 * relative to some address. Used, e.g. by shmem.
	 */
	unsigned long offset;
	bool transfer_lazy;

	/* private data for every page-xfer engine */
	union {
		struct /* local */ {
			struct cr_img *pmi; /* pagemaps */
			struct cr_img *pi;  /* pages */
		};

		struct /* page-server */ {
			int sk;
			u64 dst_id;
		};
	};

	struct page_read *parent;
};

struct cr_img {
	union {
		struct bfd _x;
		struct {
			int fd; /* should be first to coincide with _x.fd */
			int type;
			unsigned long oflags;
			char *path;
		};
	};
};

```

## 测试

```shell

make clean && make DEBUG=1
make install

cd /opt/container-migrator/workload/redis
runc run myredis
runc list

pid=13903 
criu pre-dump -vvvv --log-file log.txt --images-dir /opt/container-migrator/client_repo/gdbmyredis/checkTest -t $pid --track-mem
```