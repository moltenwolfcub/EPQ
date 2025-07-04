#include <assimp/version.h>
#include <iostream>

using namespace std;

int main()
{
	cout << aiGetVersionMajor() << "." << aiGetVersionMinor() << "." << aiGetVersionPatch() << "\n";

	return 0;
}
