Run the application on AWS EMR:

. create a jar of the application: app.jar
	This jar must not contain the Spark dependencies

For a small test:
. on the cluster, create a new Spark application step
	. Deploy mode: cluster
	. Spark-submit options: --class ETL
	. Application location: s3://mdk-phase1-q2/app.jar
	. Arguments: s3://cmucc-datasets/twitter/f15/part-00001 s3://mdk-phase1-q2/output
		(I think the folder "output" shouldn't already exist)

