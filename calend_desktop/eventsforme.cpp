#include "eventsforme.h"
#include "ui_eventsforme.h"

eventsforme::eventsforme(QString uid, QWidget *parent)
    : QWidget(parent)
    , ui(new Ui::eventsforme), uid(uid),
    client(settings.value("host").toString(), settings.value("port").toString())
{
    ui->setupUi(this);
    loadEvents();
    ui->scrollArea_eventsforme->setWidgetResizable(true);
    ui->scrollArea_eventsforme->setVerticalScrollBarPolicy(Qt::ScrollBarAsNeeded);
}

eventsforme::~eventsforme()
{
    delete ui;
}

void eventsforme::loadEvents()
{
    QWidget* eventsContainer = ui->scrollArea_eventsforme->widget();
    if (!eventsContainer) {
        eventsContainer = new QWidget();
        ui->scrollArea_eventsforme->setWidget(eventsContainer);
    }

    QLayout *oldLayout = eventsContainer->layout();
    if (oldLayout) {
        while (QLayoutItem *item = oldLayout->takeAt(0)) {
            delete item->widget();
            delete item;
        }
        delete oldLayout;
    }

    QVBoxLayout *eventsLayout = new QVBoxLayout(eventsContainer);
    eventsContainer->setLayout(eventsLayout);

    QVector<EventData> allEvents = client.getAllUserEvents(uid);

    if (allEvents.isEmpty()) {
        QLabel *noEventsLabel = new QLabel("Событий не найдено");
        noEventsLabel->setAlignment(Qt::AlignCenter);
        eventsLayout->addWidget(noEventsLabel);
        return;
    }

    bool hasEvents = false;
    for (const EventData &event : allEvents) {
        bool isCreator = (event.master == uid);
        bool isParticipant = false;
        bool isAccepted = false;

        for (const Participant &p : event.parts) {
            if (p.uid == uid) {
                isParticipant = true;
                isAccepted = p.accepted;
                break;
            }
        }

        bool showEvent = (isCreator && ui->checkBox_created->isChecked()) ||
                         (isParticipant && isAccepted && ui->checkBox_accepted->isChecked()) ||
                         (isParticipant && !isAccepted && ui->checkBox_not_accepted->isChecked());

        if (showEvent) {
            event_entry *entry = new event_entry(event, uid);
            eventsLayout->addWidget(entry);
            hasEvents = true;

            connect(entry, &event_entry::deleted, this, &eventsforme::loadEvents);
            connect(entry, &event_entry::edited, this, &eventsforme::loadEvents);
        }
    }

    if (!hasEvents) {
        QLabel *noFilteredEvents = new QLabel("Нет событий по выбранным фильтрам");
        noFilteredEvents->setAlignment(Qt::AlignCenter);
        eventsLayout->addWidget(noFilteredEvents);
    }
    ui->scrollArea_eventsforme->setWidgetResizable(true);
    eventsContainer->adjustSize();
}


void eventsforme::on_checkBox_created_stateChanged(int)
{
    loadEvents();
}


void eventsforme::on_checkBox_accepted_stateChanged(int)
{
    loadEvents();
}


void eventsforme::on_checkBox_not_accepted_stateChanged(int)
{
    loadEvents();
}

