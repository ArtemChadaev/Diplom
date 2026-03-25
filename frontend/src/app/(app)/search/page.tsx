import { Suspense } from "react";
import { SearchInterface } from "./search-interface";

export default function SearchPage() {
    return (
        <div className="flex flex-col gap-10">
            <Suspense fallback={<div className="h-64 flex items-center justify-center">Загрузка поиска...</div>}>
                <SearchInterface />
            </Suspense>
        </div>
    );
}
