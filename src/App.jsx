import { useEffect, useState } from "react";
import "./App.css";
import { DropdownButtons } from "./components";
import { team, content } from "./data";
import ReactMarkdown from 'react-markdown'
import remarkGfm from 'remark-gfm'
function App() {
  const [markdown, setMarkdown] = useState(`Just a link: https://dup.company.`)
  useEffect(()=>{
    (async ()=>{
      try {
        const res = await fetch("https://raw.githubusercontent.com/KM8Oz/WhizApp/main/docs/README.md")
        const _markdown = await res.text()
        if (_markdown.length > 0){
          setMarkdown(_markdown)
        }
      } catch (error) {
        setMarkdown("## error \n ### > loading error !")
      }
    })()
    
  },[])
  return (
    <>
      <div className="container flex flex-col mx-auto">
        <img className="p-5 rounded-full w-[300px] mx-auto animate-pulse" src="Icon.png" alt="icon" />
        <h1 className="font-extrabold p-5">{content.main.title}</h1>
        <h3 className="uppercase p-5">{content.main.description}</h3>
        <p className="uppercase p-4">{content.main.intro}</p>
        <a href={content.main.repository} className="p-3">
          {content.main.repository}
        </a>
        <DropdownButtons />
      </div>

      <main className="flex-box w-screen">
        <div className="scroll-tap">
          <h2 className="heading">Features</h2>
          <div className="flex flex-wrap w-full justify-center">
            <div className="packgroup">
              <h6 className="headline">{content.Features[0].feature}</h6>
              <p className="intro">{content.Features[0].intro}</p>
            </div>
            <div className="packgroup">
              <h6 className="headline">{content.Features[1].feature}</h6>
              <p className="intro">{content.Features[1].intro}</p>
            </div>
            <div className="packgroup">
              <h6 className="headline">{content.Features[2].feature}</h6>
              <p className="intro">{content.Features[2].intro}</p>
            </div>
            <div className="packgroup">
              <h6 className="headline">{content.Features[3].feature}</h6>
              <p className="intro">{content.Features[3].intro}</p>
            </div>
            <div className="packgroup">
              <h6 className="headline">{content.Features[4].feature}</h6>
              <p className="intro">{content.Features[4].intro}</p>
            </div>
          </div>
        </div>
        <hr />
        <div className="scroll-tap">
          <h2 className="heading">Usage</h2>
          <ReactMarkdown children={markdown} remarkPlugins={[remarkGfm]} />,
        </div>

        <hr />
        <div className="scroll-tap">
          <h2 className="heading">Configuration</h2>
          <div className="flex-box">
            <div className="grouplist">
              <h3 className="headlist">{content.Configuration[0].config}</h3>
              <p className="mainlist">{content.Configuration[0].intro}</p>
            </div>
            <div className="grouplist">
              <h3 className="headlist">{content.Configuration[1].config}</h3>
              <p className="mainlist">{content.Configuration[1].intro}</p>
            </div>
            <div className="grouplist">
              <h3 className="headlist">{content.Configuration[2].config}</h3>
              <p className="mainlist">{content.Configuration[2].intro}</p>
            </div>
            <div className="grouplist">
              <h3 className="headlist">{content.Configuration[3].config}</h3>
              <p className="mainlist">{content.Configuration[3].intro}</p>
            </div>
          </div>
        </div>

        
<div className="container scroll-tap my-24 mx-auto md:px-6">
  <section className="scroll-tap mb-32 text-center">
    <h2 className="heading mb-20 text-3xl font-bold">Acknowledgements</h2>

    <div className="grid lg:grid-cols-3 lg:gap-x-12 sm:grid-cols-1 px-10">
      <div className="mb-12 lg:mb-0">
        <div className="block h-full rounded-lg bg-neutral-900 shadow-[0_2px_15px_-3px_rgba(0,0,0,0.07),0_10px_20px_-2px_rgba(0,0,0,0.04)] ">
          <div className="flex justify-center">
            <div className="-mt-8 inline-block rounded-full bg-primary-100 p-4 text-primary shadow-md">
            <img
              className="team-img"
              src='gpt.jpeg'
              alt="team"
            />
            </div>
          </div>
          <div className="p-6">
            <h5 className="mb-4 text-lg font-semibold">{content.Acknowledgements[0].skill}</h5>
            <p>
              {content.Acknowledgements[0].description}
            </p>
          </div>
        </div>
      </div>

      <div className="mb-12 lg:mb-0">
        <div className="block h-full rounded-lg bg-neutral-900 shadow-[0_2px_15px_-3px_rgba(0,0,0,0.07),0_10px_20px_-2px_rgba(0,0,0,0.04)] ">
          <div className="flex justify-center">
            <div className="-mt-8 inline-block rounded-full bg-primary-100 p-4 text-primary shadow-md">
            <img
              className="team-img"
              src='golang.jpeg'
              alt="team"
            />
            </div>
          </div>
          <div className="p-6">
            <h5 className="mb-4 text-lg font-semibold">{content.Acknowledgements[1].skill}</h5>
            <p>
             {content.Acknowledgements[1].description}
            </p>
          </div>
        </div>
      </div>

      <div className="">
        <div className="block h-full rounded-lg bg-neutral-900 shadow-[0_2px_15px_-3px_rgba(0,0,0,0.07),0_10px_20px_-2px_rgba(0,0,0,0.04)] bg-neutral-900">
          <div className="flex justify-center">
            <div className="-mt-8 inline-block rounded-full bg-primary-100 p-4 text-primary shadow-md">
            <img
              className="team-img"
              src="github.jpeg"
              alt="team"
            />
            </div>
          </div>
          <div className="p-6">
            <h5 className="mb-4 text-lg font-semibold">{content.Acknowledgements[2].skill}</h5>
            <p>
           {content.Acknowledgements[2].description}
            </p>
          </div>
        </div>
      </div>
    </div>
  </section>

</div>

      </main>

      <footer className="flex w-full justify-center scroll-tap mx-auto mb-0">
        <p className="text-4xl font-semibold">
          Team 
        </p>
        <div className="flex flex-row flex-wrap justify-between px-3">
          <div className="colab">
            <img
              className="team-img"
              src={team.contributor1.image}
              alt="team"
            />
            <h4 className="p-5">{team.contributor1.name}</h4>
            <p>{team.contributor1.job}</p>
          </div>

          <div className="colab">
            <img
              className="team-img"
              src={team.contributor2.image}
              alt="team"
            />
            <h4 className="p-5">{team.contributor2.name}</h4>
            <p>{team.contributor2.job}</p>
          </div>
          <div className="colab">
            <img
              className="team-img"
              src={team.contributor3.image}
              alt="team"
            />
            <h4 className="p-5">{team.contributor3.name}</h4>
            <p>{team.contributor3.job}</p>
          </div>
          <div className="colab">
            <img
              className="team-img"
              src={team.contributor4.image}
              alt="team"
            />
            <h4 className="p-5">{team.contributor4.name}</h4>
            <p>{team.contributor4.job}</p>
          </div>
          <div className="colab">
            <img
              className="team-img"
              src={team.contributor5.image}
              alt="team"
            />
            <h4 className="p-5">{team.contributor5.name}</h4>
            <p>{team.contributor5.job}</p>
          </div>
        </div>
        <div className="bg-neutral-200 p-6 text-center dark:bg-neutral-700 flex-end w-screen">
          <span className="">© 2023 By💟. All Right Reserved.</span>
          <a
            className="font-semibold text-neutral-600 dark:text-neutral-400"
            href={content.main.link}
          >
            WhizApp
          </a>
        </div>
      </footer>
    </>
  );
}

export default App;
