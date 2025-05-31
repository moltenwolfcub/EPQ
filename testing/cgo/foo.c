#include <stdio.h>

int getSchoolYear(int age)
{
	if (age < 5)
	{
		return 100;
	}
	else
	{
		return (age - 5);
	}
}
