#include <stdio.h>

int getSchoolYear(int age)
{
	if (age < 5)
	{
		return 0;
	}
	else
	{
		return (age - 5);
	}
}
