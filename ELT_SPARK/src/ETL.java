import com.google.gson.Gson;
import org.apache.log4j.Level;
import org.apache.log4j.Logger;
import org.apache.spark.SparkConf;
import org.apache.spark.api.java.JavaRDD;
import org.apache.spark.api.java.JavaSparkContext;
import org.apache.spark.api.java.function.Function;

import java.io.BufferedReader;
import java.io.FileInputStream;
import java.io.FileNotFoundException;
import java.io.FileReader;

/**
 * Extract - transform - load using Spark
 */
public class ETL {

    public static void main(String[] args) throws Exception {

        Logger.getLogger("org.apache.spark").setLevel(Level.ERROR);
        Logger.getLogger("org.apache.spark.storage.BlockManager").setLevel(Level.ERROR);

        if (args.length != 2)
            throw new IllegalArgumentException("Two arguments expected");

        String inputLocation = args[0];  // s3 folder/file
        String outputLocation = args[1];

        SparkConf conf = new SparkConf().setAppName("ETL");
        JavaSparkContext sc = new JavaSparkContext(conf);

        JavaRDD<String> lines = sc.textFile(inputLocation);
        JavaRDD<Tweet> tweets = lines.map(new Function<String, Tweet>() {
            public Tweet call(String s) {
                Gson gson = new Gson();
                Tweet tweet = null;
                try {
                    tweet = gson.fromJson(s, Tweet.class);
                } catch (Exception ignored) {}
                return tweet;
            }
        });

        tweets.map(new Function<Tweet, Long>() {
            public Long call(Tweet tweet) {
                return tweet.id;
            }
        }).saveAsTextFile(outputLocation);

        /*
        Gson gson = new Gson();
        Tweet tweet;
        BufferedReader br = new BufferedReader(new FileReader("/home/dg/cmu/mdk/elt_spark/src/tt.txt"));
        String line;
        while ((line = br.readLine()) != null) {
            try {
                tweet = gson.fromJson(line, Tweet.class);
                System.out.println("id: " + tweet.id);
                System.out.println("user: " + tweet.user.id_str);
            } catch (Exception ignored) {}
        }
        */
    }

}
