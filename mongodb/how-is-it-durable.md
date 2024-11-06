# How it maintain durability
MongoDB will use WAL to write data to a on-disk journal, so that in the event of crash, MongoDB can use this journal to apply the write operations to the data files.

Journal is a sequential, binary transaction log used to bring the db into a valid state in the event of a hard shutdown.

Journaling writes data first to the journal and then to the core data files. MongoDB enables journaling by default for 64-bit builds of MongoDB version 2.0 and newer. Journal files are pre-allocated and exist as files in the data directory.

The operation that "flush" the data from journal files to the data files is done lazily. (about 60 sec)