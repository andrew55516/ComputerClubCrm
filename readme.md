# ComputerClubCrm

## Start
___
Download the repository and run following commands from the root project directory:
```
$ docker build -t computerclubcrm:latest .
```
```
$ docker run computerclubcrm test_in.txt
```
You can replace test_in.txt with test1_in.txt or test2_in.txt, and you also can add your own test files.

## Details
___
I added more errors for the incoming event ID3 (client is waiting):
- If the client has not already come - 'ClientUnknown' error
- If the client has already sat down - 'ClientHasAlreadySatDown' error
- If the client has already waited - 'ClientHasAlreadyWaited' error