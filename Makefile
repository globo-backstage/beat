setup:
	@pip install -e .\[tests\]

run:
	@pserve development.ini --reload
