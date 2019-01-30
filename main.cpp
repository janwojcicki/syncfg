#include <iostream>
#include <cstring>
#include <unistd.h>
#include <sys/types.h>
#include <string>
#include <fstream>
#include <pwd.h>
#include <vector>
struct passwd *pw = getpwuid(getuid());

const char *homedir = pw->pw_dir;
const std::string hdir (homedir);

int main(int argc, char* argv[]){
	std::vector<std::string> config;
	std::ifstream inp (hdir+"/.kr");
	std::string line;
	while(std::getline(inp, line)){
		config.push_back(line);
	}
	inp.close();
	if(argc > 1){
		if (strcmp(argv[1], "add") && argc == 4){
			for (std::string s : config){
				if (s.substr(0,1) == "#"){
					if (s.substr(1, s.size()-1) == argv[2])	{
						std::cout << "FOUND";
					}
				}
			}
		}	
	}	
	return 0;
}
