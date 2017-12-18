#!/usr/bin/python
import sys
import json

def main():
    """
    main()
    Allows for json stdin, then stdout the response.
    """

    for line in sys.stdin:
        line = json.loads(line.strip("\n"))
        line += "-from-plugin"
        print json.dumps(line)

main()
