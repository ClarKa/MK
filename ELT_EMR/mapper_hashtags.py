#!/usr/bin/env python

import sys
import os
import json
import codecs
import datetime
import re
from datetime import timedelta

def main(argv):
	r = re.compile(r'\w+')
	line = sys.stdin.readline()
	try:
		while line:			
			tw = json.loads(line)
			text = tw['text']
			date_text = tw['created_at']
			timezone = date_text.split()[-2]
			date = datetime.datetime.strptime(date_text, "%a %b %d %H:%M:%S "+timezone+" %Y")
			pos = 0
			if timezone[0] == "+":
				pos = 1
			elif timezone[0] == "-":
				pos = -1
			hr = int(timezone[1:3])
			mn = int(timezone[3:5])
			date = date+timedelta(hours=pos*hr,minutes=pos*mn)
			hashtags = tw['entities']['hashtags']
			for hashtag in hashtags:
				l = json.dumps(hashtag["text"])[1:-1]+'\t'+str(tw['id'])+'\t'+tw["user"]["id_str"]+'\t'+str(date)
				print l
			line =  sys.stdin.readline()
	except "end of file":
		return None
if __name__ == "__main__":
	main(sys.argv)