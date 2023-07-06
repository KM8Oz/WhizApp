#!/usr/bin/python3
import re
with open("./dataset_darija_stories.txt", 'r',  encoding="utf-8") as file:
    sentences = file.readlines()
    cleaned_sentences = ["{}\n".format(re.sub(
        '[\(\)\{\}<>!(**)*;.\n]', '', el).__str__().strip()) for el in sentences if el != '\n']
    print(str.join("", cleaned_sentences)[0:100])
with open("./dataset_darija_stories.txt", 'w', encoding="utf-8") as file:
    file.writelines(cleaned_sentences)
