export type Note = {
    id: string;
    title: string;
    content: string;
    createdAt: Date;
    priority?: string;
    completionTime?: Date;
};
